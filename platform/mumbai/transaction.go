package mumbai

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/address"
	"math/big"
	"strings"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

var ErrTransferActionUnmarshal = errors.New("unable marshaling to transfer action struct")

const hashEventTypeTransfer = "0xddf252ad"

func (p *Platform) NormalizeRawBlock(rawBlock []byte) (*types.Block, error) {
	var block Block
	if err := json.Unmarshal(rawBlock, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return p.NormalizeBlock(block, block.TxReceipts.Map()), nil
}

func (p *Platform) NormalizeBlock(block Block, receipts map[string]TransactionReceipt) *types.Block {
	size, _ := new(big.Float).SetInt((*big.Int)(block.Size)).Float64()

	uncleHash := ""
	for _, uncle := range block.Uncles {
		uncleHash = uncleHash + "," + uncle.(string)
	}

	normalizeBlock := types.Block{
		Hash:                block.Hash,
		Number:              (*big.Int)(block.Number).Uint64(),
		Time:                (*big.Int)(block.Timestamp).Uint64(),
		ParentHash:          block.ParentHash,
		Difficulty:          (*big.Int)(block.Difficulty).String(),
		GasUsed:             (*big.Int)(block.GasUsed).Uint64(),
		GasLimit:            (*big.Int)(block.GasLimit).Uint64(),
		Nonce:               (*big.Int)(block.Nonce).String(),
		Miner:               block.Miner,
		Size:                size,
		StateRootHash:       block.StateRoot,
		UncleHash:           uncleHash,
		TransactionRootHash: block.TransactionsRoot,
		ReceiptRootHash:     block.ReceiptsRoot,
		ExtraData:           []byte(block.ExtraData),
	}
	txs := make(types.Txs, 0)

	for _, tx := range block.Transactions {
		if tx.BlockNumber == nil { // pending transaction
			continue
		}
		receipt, ok := receipts[tx.Hash]
		if !ok {
			continue
		}
		if receipt.Status == nil {
			log.WithField("chain", p.Coin().Handle).Errorf("Empty tx receipt status for tx: %s", tx.Hash)

			continue
		}
		normalized := p.NormalizeTransaction(receipt, tx, block)
		if normalized == nil {
			continue
		}

		status := types.StatusSuccess
		if (*big.Int)(receipt.Status).Cmp(big.NewInt(0)) == 0 {
			status = types.StatusError
		}
		normalized.Status = status

		txs = append(txs, *normalized)
	}
	normalizeBlock.Txs = txs

	return &normalizeBlock
}

func (p *Platform) NormalizeTransaction(receipt TransactionReceipt, tx Transaction, block Block) *types.Tx {
	ts := (*big.Int)(block.Timestamp).Int64()

	if tx.Input == "0x" { // coin transfer
		return p.NormalizeCoinTransfer(tx, receipt, ts, block.BaseFeePerGas)
	}

	isTransfer, transferLog := p.checkTransfer(receipt.Logs)
	if isTransfer {
		return p.NormalizeTokenTransfer(tx, receipt, ts, transferLog, block.BaseFeePerGas)
	}

	return p.NormalizeContractCall(tx, receipt, ts, block.BaseFeePerGas)
}

func (p *Platform) normalizeBaseOfTx(srcTx Transaction,
	receipt TransactionReceipt, timestamp int64, baseFeePerGas *types.HexNumber,
) *types.Tx {
	l := log.WithFields(log.Fields{"chain": p.Coin().Handle, "tx": srcTx.Hash})

	addressFrom, err := address.EIP55Checksum(srcTx.From)
	if err != nil {
		l.WithError(err).Error("Error while getting eip55 checksum of ethereum address 'from'")

		return nil
	}

	addressTo, err := address.EIP55Checksum(srcTx.To)
	if err != nil {
		l.WithError(err).Error("Error while getting eip55 checksum of ethereum address 'to'")

		return nil
	}

	events := p.NormalizeEventLog(receipt.Logs)

	tx := &types.Tx{
		Hash:           srcTx.Hash,
		Chain:          p.Coin().Handle,
		From:           addressFrom,
		To:             addressTo,
		BlockCreatedAt: timestamp,
		Block:          (*big.Int)(srcTx.BlockNumber).Uint64(),
		Fee: types.Fee{
			Asset: p.Coin().AssetID(),
		},
		BlockHash: srcTx.BlockHash,
		Events:    events,
	}

	if srcTx.Nonce == nil {
		l.Errorf("nonce is nil")

		return nil
	}
	tx.Sequence = (*big.Int)(srcTx.Nonce).Uint64()

	feeValue, err := srcTx.Fee((*big.Int)(baseFeePerGas), (*big.Int)(receipt.GasUsed))
	if err != nil {
		l.WithError(err).Error("cannot compute fee")

		return nil
	}

	tx.Fee.Amount = types.Amount(feeValue)

	return tx
}

func (p *Platform) NormalizeEventLog(eventLogs []EventLog) []types.Event {
	events := make([]types.Event, 0)
	for _, eventLog := range eventLogs {
		events = append(events, types.Event{
			Address:          eventLog.Address,
			Topics:           eventLog.Topics,
			Data:             eventLog.Data,
			BlockNumber:      eventLog.BlockNumber,
			TransactionHash:  eventLog.TransactionHash,
			TransactionIndex: eventLog.TransactionIndex,
			BlockHash:        eventLog.BlockHash,
			LogIndex:         eventLog.LogIndex,
			Removed:          eventLog.Removed,
		})
	}

	return events
}

func (p *Platform) NormalizeCoinTransfer(srcTx Transaction,
	receipt TransactionReceipt, timestamp int64, baseFeePerGas *types.HexNumber,
) *types.Tx {
	tx := p.normalizeBaseOfTx(srcTx, receipt, timestamp, baseFeePerGas)
	if tx == nil {
		return nil
	}
	tx.Type = types.TxTransfer

	if srcTx.Value == nil {
		log.WithField("tx", srcTx.Hash).Error("transfer value is nil")

		return nil
	}
	tx.Metadata = &types.Transfer{
		Asset:  p.Coin().AssetID(),
		Amount: types.Amount((*big.Int)(srcTx.Value).String()),
	}

	return tx
}

func (p *Platform) NormalizeTokenTransfer(srcTx Transaction, receipt TransactionReceipt,
	timestamp int64, eventLog EventLog, baseFeePerGas *types.HexNumber,
) *types.Tx {
	tx := p.normalizeBaseOfTx(srcTx, receipt, timestamp, baseFeePerGas)
	if tx == nil {
		return nil
	}

	l := log.WithFields(log.Fields{"chain": p.Coin().Handle, "tx": tx.Hash})

	addressFrom, err := address.EIP55Checksum(common.HexToAddress(eventLog.Topics[1]).Hex())
	if err != nil {
		l.WithError(err).Errorf("Could not get eip55 checksum of ethereum token sender")

		return nil
	}

	addressTo, err := address.EIP55Checksum(common.HexToAddress(eventLog.Topics[2]).Hex())
	if err != nil {
		l.WithError(err).Errorf("Could not get eip55 checksum of ethereum token recipient")

		return nil
	}

	value, err := hexDataToBigInt(eventLog.Data)
	if err != nil {
		l.WithError(err).Error("Hex conversion to big int error")

		return nil
	}

	token, err := address.EIP55Checksum(eventLog.Address) // Use address of a contract that was transferred
	if err != nil {
		l.WithError(err).Error("EIP55 checksum getting error")

		return nil
	}

	tx.From = addressFrom
	tx.To = addressTo
	tx.Type = types.TxTransfer
	tx.Metadata = &types.Transfer{
		Asset:  p.Coin().TokenAssetID(token), // Token contract
		Amount: types.Amount(value.String()),
	}

	return tx
}

func (p *Platform) NormalizeContractCall(srcTx Transaction,
	receipt TransactionReceipt, timestamp int64, baseFeePerGas *types.HexNumber,
) *types.Tx {
	tx := p.normalizeBaseOfTx(srcTx, receipt, timestamp, baseFeePerGas)
	if tx == nil {
		return nil
	}

	tx.Type = types.TxContractCall

	if srcTx.Value == nil {
		log.WithField("tx", srcTx.Hash).Error("contract call value is nil")

		return nil
	}
	tx.Metadata = &types.ContractCall{
		Asset:  p.Coin().AssetID(),
		Amount: types.Amount((*big.Int)(srcTx.Value).String()),
		Input:  srcTx.Input,
	}

	return tx
}

//nolint:goerr113
func hexDataToBigInt(data string) (*big.Int, error) {
	data = strings.Replace(data, "0x", "", 1)

	value, ok := new(big.Int).SetString(data, 16)
	if !ok {
		return nil, fmt.Errorf("cannot convert data to *big.Int: %s", data)
	}

	return value, nil
}

func (p *Platform) checkTransfer(logs []EventLog) (bool, EventLog) {
	var transferLog EventLog
	transferLogsNum := 0
	for i, eventLog := range logs {
		topics := eventLog.Topics
		if len(topics) == 3 && strings.HasPrefix(topics[0], hashEventTypeTransfer) {
			transferLog = logs[i]
			transferLogsNum++
		}
	}

	return transferLogsNum == 1, transferLog
}

func (p *Platform) DecodeTransactionLogs(logs []EventLog) []LogTransfer {
	var transferEvents []LogTransfer
	var transferEvent LogTransfer
	transferEventHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	for _, vLog := range logs {
		if strings.Compare(vLog.Topics[0], transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				transferEvent.From = common.HexToAddress(vLog.Topics[1])
				transferEvent.To = common.HexToAddress(vLog.Topics[2])
				transferEvent.TokenId = common.HexToHash(vLog.Topics[3]).Big()
				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	return transferEvents
}
