package types

type (
	Amount          string
	Asset           string
	Status          string
	TransactionType string
)

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type (
	Block struct {
		Number int64 `json:"number"`
		Txs    []Tx  `json:"txs"`
	}

	Txs []Tx

	Tx struct {
		Hash           string          `json:"hash" bson:"hash"`
		From           string          `json:"from" bson:"from"`
		To             string          `json:"to" bson:"to"`
		BlockCreatedAt int64           `json:"block_created_at" bson:"block_created_at"`
		Block          uint64          `json:"block_num" bson:"block_num"`
		Sequence       uint64          `json:"sequence" bson:"status"`
		Status         Status          `json:"status" bson:"status"`
		Type           TransactionType `json:"type" bson:"type"`
		Fee            Fee             `json:"fee" bson:"type"`
		Metadata       interface{}     `json:"metadata" bson:"metadata"`
	}

	Fee struct {
		Asset  Asset  `json:"asset"`
		Amount Amount `json:"amount"`
	}
)
