package types

import (
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
)

type (
	Amount          string
	Asset           string
	Status          string
	TransactionType string
)

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"

	TxTransfer          TransactionType = "transfer"
	TxContractCall      TransactionType = "contract_call"
	TxStakeClaimRewards TransactionType = "stake_claim_rewards"
	TxStakeDelegate     TransactionType = "stake_delegate"
	TxStakeUndelegate   TransactionType = "stake_undelegate"
	TxStakeRedelegate   TransactionType = "stake_redelegate"
)

type (
	Block struct {
		Hash                string  `json:"string"`
		Number              uint64  `json:"number"`
		Time                uint64  `json:"time"`
		ParentHash          string  `json:"parent_hash"`
		Difficulty          string  `json:"difficulty"`
		GasUsed             uint64  `json:"gas_used"`
		GasLimit            uint64  `json:"gas_limit"`
		Nonce               string  `json:"nonce"`
		Miner               string  `json:"miner"`
		Size                float64 `json:"size"`
		StateRootHash       string  `json:"state_root_hash"`
		UncleHash           string  `json:"uncle_hash"`
		TransactionRootHash string  `json:"tx_root_hash"`
		ReceiptRootHash     string  `json:"receipt_root_hash"`
		ExtraData           []byte  `json:"extra_data"`
		Txs                 []Tx    `json:"txs"`
	}

	Txs []Tx

	Tx struct {
		Hash           string          `json:"hash" bson:"hash"`
		Chain          string          `json:"chain" bson:"chain"`
		From           string          `json:"from" bson:"from"`
		To             string          `json:"to" bson:"to"`
		BlockCreatedAt int64           `json:"block_created_at" bson:"block_created_at"`
		Block          uint64          `json:"block" bson:"block"`
		Sequence       uint64          `json:"sequence" bson:"status"`
		Status         Status          `json:"status" bson:"status"`
		Type           TransactionType `json:"type" bson:"type"`
		Memo           string          `json:"memo,omitempty"`
		Fee            Fee             `json:"fee" bson:"type"`
		Metadata       interface{}     `json:"metadata" bson:"metadata"`
		BlockHash      string          `json:"blockhash" bson:"blockhash"`
		Events         []Event         `json:"event" bson:"type"`
	}

	Fee struct {
		Asset  coin.AssetID `json:"asset"`
		Amount Amount       `json:"amount"`
	}

	Event struct {
		Address          string     `json:"address"`
		Topics           []string   `json:"topics"`
		Data             string     `json:"data"`
		BlockNumber      *HexNumber `json:"blockNumber"`
		TransactionHash  string     `json:"transactionHash"`
		TransactionIndex *HexNumber `json:"transactionIndex"`
		BlockHash        string     `json:"blockHash"`
		LogIndex         *HexNumber `json:"logIndex"`
		Removed          bool       `json:"removed"`
	}
)

// Tx metadata types
type (
	Transfer struct {
		Asset  coin.AssetID `json:"asset"`
		Amount Amount       `json:"amount"`
	}

	TransferLog struct {
		From     string
		To       string
		Contract string
		TokenId  *HexNumber
		Amount   *HexNumber
	}

	ContractCall struct {
		Asset  coin.AssetID `json:"asset"`
		Amount Amount       `json:"amount"`
		Input  string       `json:"input"`
	}
)
