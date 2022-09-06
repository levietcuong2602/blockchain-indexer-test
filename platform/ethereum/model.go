package ethereum

import "github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"

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
