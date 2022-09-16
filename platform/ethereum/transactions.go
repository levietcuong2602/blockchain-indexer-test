package ethereum

// import (
// 	"encoding/json"
// 	"math/big"

// 	log "github.com/sirupsen/logrus"

// 	"github.com/unanoc/blockchain-indexer/pkg/primitives/address"
// 	"github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"
// )

// func (p *Platform) NormalizeRawBlock(rawBlock []byte) (types.Txs, error) {
// 	var block Block
// 	if err := json.Unmarshal(rawBlock, &block); err != nil {
// 		return nil, err
// 	}

// 	return p.NormalizeBlock(block, block.TxReceipts.Map()), nil
// }

// func (p *Platform) NormalizeBlock(block Block, receipts map[string]TransactionReceipt) types.Txs {
// 	var txs types.Txs

// 	for _, tx := range block.Transactions {
// 		if tx.BlockNumber == nil { // pending transaction
// 			continue
// 		}

// 		receipt, ok := receipts[tx.Hash]
// 		if !ok {
// 			continue
// 		}

// 		if receipt.Status == nil {
// 			log.WithField("chain", p.chain).Errorf("Empty tx receipt status for tx: %s", tx.Hash)

// 			continue
// 		}

// 		normalized := p.NormalizeTransaction(receipt, &tx, block)
// 		if normalized == nil {
// 			continue
// 		}

// 		status := types.StatusSuccess
// 		if (*big.Int)(receipt.Status).Cmp(big.NewInt(0)) == 0 {
// 			status = types.StatusError
// 		}
// 		normalized.Status = status

// 		txs = append(txs, *normalized)
// 	}

// 	return txs
// }

// func (p *Platform) NormalizeTransaction(receipt TransactionReceipt, tx *Transaction, block Block) *types.Tx {
// 	ts := (*big.Int)(block.Timestamp).Int64()

// 	if tx.Input == "0x" { // coin transfer
// 		return p.normalizeTransfer(tx, receipt, ts, block.BaseFeePerGas)
// 	}

// 	isTransfer, transferLog := p.checkTransfer(receipt.Logs)
// 	if isTransfer {
// 		return p.normalizeTokenTransfer(tx, receipt, ts, transferLog, block.BaseFeePerGas)
// 	}

// 	isSwap, transferFromLog, transferToLog := p.checkSwap(receipt.Logs, tx.From, tx.To)
// 	if isSwap {
// 		normalizedTx := p.normalizeSwap(tx, receipt, ts, transferFromLog, transferToLog, block.BaseFeePerGas)
// 		if normalizedTx != nil {
// 			return normalizedTx
// 		}
// 	}

// 	return p.normalizeContractCall(tx, receipt, ts, block.BaseFeePerGas)
// }

// func (ec *Client) normalizeBaseOfTx(srcTx *Transaction,
// 	receipt TransactionReceipt, timestamp int64, baseFeePerGas *types.HexNumber,
// ) *types.Tx {
// 	logger := log.WithFields(log.Fields{"handle": ec.coin.Handle, "tx": srcTx.Hash})

// 	addressFrom, err := address.EIP55Checksum(srcTx.From)
// 	if err != nil {
// 		logger.Errorf("Could not get eip55 checksum of ethereum address from: %v", err)
// 		return nil
// 	}

// 	addressTo, err := address.EIP55Checksum(srcTx.To)
// 	if err != nil {
// 		logger.Errorf("Could not get eip55 checksum of ethereum address to: %v", err)
// 		return nil
// 	}

// 	tx := &types.Tx{
// 		Hash:           srcTx.Hash,
// 		From:           addressFrom,
// 		To:             addressTo,
// 		BlockCreatedAt: timestamp,
// 		Block:          (*big.Int)(srcTx.BlockNumber).Uint64(),
// 		Fee: types.Fee{
// 			Asset: ec.coin.AssetID(),
// 		},
// 	}

// 	if srcTx.Nonce == nil {
// 		logger.Errorf("nonce is nil")
// 		return nil
// 	}
// 	tx.Sequence = (*big.Int)(srcTx.Nonce).Uint64()

// 	if feeValue, err := srcTx.Fee((*big.Int)(baseFeePerGas), (*big.Int)(receipt.GasUsed)); err != nil {
// 		logger.WithError(err).Error("cannot compute fee")
// 		return nil
// 	} else {
// 		tx.Fee.Value = types.Amount(feeValue)
// 	}

// 	return tx
// }
