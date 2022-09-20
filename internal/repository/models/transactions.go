package models

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

type Transaction struct {
	Hash           string `gorm:"primary_key"`
	Chain          string `gorm:"not_null"`
	Sender         string `gorm:"type:varchar(256); not_null"`
	Recipient      string `gorm:"type:varchar(256); not_null"`
	Fee            string
	Block          uint64
	BlockCreatedAt int64
	Sequence       uint64
	Status         types.Status
	Type           types.TransactionType
	Metadata       postgres.Jsonb
	CreatedAt      int64 `gorm:"auto_create_time"`
}

func NormalizeTransaction(tx types.Tx, chain string) (*Transaction, error) {
	transaction := Transaction{
		Hash:           tx.Hash,
		Chain:          chain,
		Sender:         tx.From,
		Recipient:      tx.To,
		Block:          tx.Block,
		BlockCreatedAt: tx.BlockCreatedAt,
		Sequence:       tx.Sequence,
		Status:         tx.Status,
		Type:           tx.Type,
		Fee:            string(tx.Fee.Amount),
	}

	metadataRaw, err := json.Marshal(tx.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	transaction.Metadata = postgres.Jsonb{RawMessage: metadataRaw}

	return &transaction, nil
}

func NormalizeTransactions(txs types.Txs, chain string) ([]Transaction, error) {
	result := make([]Transaction, len(txs))
	for i := range txs {
		normalizedTx, err := NormalizeTransaction(txs[i], chain)
		if err != nil {
			return nil, err
		}

		result[i] = *normalizedTx
	}

	return result, nil
}
