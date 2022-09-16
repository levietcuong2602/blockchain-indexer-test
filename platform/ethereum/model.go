package ethereum

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

const txTypeEIP1559 = 2

// Block By Number
type (
	Block struct {
		Hash          string              `json:"hash"`
		Number        *types.HexNumber    `json:"number"`
		Timestamp     *types.HexNumber    `json:"timestamp"`
		ParentHash    string              `json:"parentHash"`
		Transactions  []Transaction       `json:"transactions"`
		BaseFeePerGas *types.HexNumber    `json:"baseFeePerGas"`
		TxReceipts    TransactionReceipts `json:"tx_receipts"`
	}

	Transaction struct {
		BlockNumber          *types.HexNumber `json:"blockNumber,omitempty"`
		Hash                 string           `json:"hash"`
		From                 string           `json:"from"`
		To                   string           `json:"to"`
		Value                *types.HexNumber `json:"value"`
		Gas                  *types.HexNumber `json:"gas"`
		GasPrice             *types.HexNumber `json:"gasPrice"`
		Input                string           `json:"input"`
		Index                *types.HexNumber `json:"transactionIndex"`
		Nonce                *types.HexNumber `json:"nonce"`
		MaxPriorityFeePerGas *types.HexNumber `json:"maxPriorityFeePerGas"`
		Type                 *types.HexNumber `json:"type"`
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

	TransactionReceipts []TransactionReceipt
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
