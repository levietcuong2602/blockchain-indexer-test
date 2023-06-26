package mumbai

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
	"math/big"
)

const txTypeEIP1559 = 2

// Block By Number
type (
	Block struct {
		BaseFeePerGas    *types.HexNumber    `json:"baseFeePerGas"`
		Difficulty       *types.HexNumber    `json:"difficulty"`
		ExtraData        string              `json:"extraData"`
		GasLimit         *types.HexNumber    `json:"gasLimit"`
		GasUsed          *types.HexNumber    `json:"gasUsed"`
		Hash             string              `json:"hash"`
		LogsBloom        string              `json:"logsBloom"`
		Miner            string              `json:"miner"`
		MixHash          string              `json:"mixHash"`
		Nonce            *types.HexNumber    `json:"nonce"`
		Number           *types.HexNumber    `json:"number"`
		ParentHash       string              `json:"parentHash"`
		ReceiptsRoot     string              `json:"receiptsRoot"`
		Sha3Uncles       string              `json:"sha3Uncles"`
		Size             *types.HexNumber    `json:"size"`
		StateRoot        string              `json:"stateRoot"`
		Timestamp        *types.HexNumber    `json:"timestamp"`
		TotalDifficulty  *types.HexNumber    `json:"totalDifficulty"`
		Transactions     []Transaction       `json:"transactions"`
		TransactionsRoot string              `json:"transactionsRoot"`
		Uncles           []interface{}       `json:"uncles"`
		TxReceipts       TransactionReceipts `json:"txReceipts"`
	}

	TransactionReceipt struct {
		BlockHash         string           `json:"blockHash"`
		BlockNumber       *types.HexNumber `json:"blockNumber"`
		ContractAddress   *string          `json:"contractAddress"`
		CumulativeGasUsed *types.HexNumber `json:"cumulativeGasUsed"`
		EffectiveGasPrice *types.HexNumber `json:"effectiveGasPrice"`
		From              string           `json:"from"`
		GasUsed           *types.HexNumber `json:"gasUsed"`
		Logs              []EventLog       `json:"logs"`
		To                string           `json:"to"`
		TransactionHash   string           `json:"transactionHash"`
		TransactionIndex  *types.HexNumber `json:"transactionIndex"`
		Type              *types.HexNumber `json:"type"`
		Status            *types.HexNumber `json:"status"`
	}

	EventLog struct {
		Address          string           `json:"address"`
		Topics           []string         `json:"topics"`
		Data             string           `json:"data"`
		BlockNumber      *types.HexNumber `json:"blockNumber"`
		TransactionHash  string           `json:"transactionHash"`
		TransactionIndex *types.HexNumber `json:"transactionIndex"`
		BlockHash        string           `json:"blockHash"`
		LogIndex         *types.HexNumber `json:"logIndex"`
		Removed          bool             `json:"removed"`
	}

	Transaction struct {
		AccessList           []interface{}    `json:"accessList,omitempty"`
		BlockHash            string           `json:"blockHash"`
		BlockNumber          *types.HexNumber `json:"blockNumber"`
		ChainID              *types.HexNumber `json:"chainId,omitempty"`
		From                 string           `json:"from"`
		Gas                  *types.HexNumber `json:"gas"`
		GasPrice             *types.HexNumber `json:"gasPrice"`
		Hash                 string           `json:"hash"`
		Input                string           `json:"input"`
		MaxFeePerGas         *types.HexNumber `json:"maxFeePerGas,omitempty"`
		MaxPriorityFeePerGas *types.HexNumber `json:"maxPriorityFeePerGas,omitempty"`
		Nonce                *types.HexNumber `json:"nonce"`
		R                    string           `json:"r"`
		S                    string           `json:"s"`
		To                   string           `json:"to"`
		TransactionIndex     *types.HexNumber `json:"transactionIndex"`
		Type                 *types.HexNumber `json:"type"`
		V                    *types.HexNumber `json:"v"`
		Value                *types.HexNumber `json:"value"`
	}

	TransactionReceipts []TransactionReceipt

	//	Event struct
	LogTransfer struct {
		From    common.Address
		To      common.Address
		TokenId *big.Int
	}
)

func (tr TransactionReceipts) Map() map[string]TransactionReceipt {
	result := make(map[string]TransactionReceipt)
	for _, receipt := range tr {
		result[receipt.TransactionHash] = receipt
	}

	return result
}

type (
	TxStatus struct {
		Message string `json:"message"`
		Result  Result `json:"result"`
	}

	Result struct {
		Status string `json:"status"`
	}
)

//nolint:goerr113
func (tx *Transaction) Fee(baseFeePerGas, gasUsed *big.Int) (string, error) {
	var txType int64
	if tx.Type != nil {
		txType = (*big.Int)(tx.Type).Int64()
	}

	switch txType {
	case txTypeEIP1559:
		if baseFeePerGas == nil {
			return "", fmt.Errorf("base fee per gas is empty for tx %s", tx.Hash)
		}
		if tx.MaxPriorityFeePerGas == nil {
			return "", fmt.Errorf("max prio fee per gas is empty for tx %s", tx.Hash)
		}
		maxPrioFeePerGas := (*big.Int)(tx.MaxPriorityFeePerGas)
		if tx.Gas == nil {
			return "", fmt.Errorf("expected gas used is empty for tx %s", tx.Hash)
		}

		tmp := &big.Int{}
		tmp.Add(baseFeePerGas, maxPrioFeePerGas)

		fee := &big.Int{}
		fee.Mul(tmp, gasUsed)

		return fee.String(), nil
	default:
		if tx.Gas == nil || tx.GasPrice == nil {
			return "", errors.New("gas and gasPrice should not be nil")
		}
		gasPrice := (*big.Int)(tx.GasPrice)

		return gasUsed.Mul(gasUsed, gasPrice).String(), nil
	}
}
