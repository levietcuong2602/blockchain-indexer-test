package models

import "github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"

type Transaction struct {
	Hash           string          `gorm:"primary_key"`
	Chain          types.ChainType `gorm:"not_null"`
	Sender         string          `gorm:"type:varchar(256); not_null"`
	Recipient      string          `gorm:"type:varchar(256); not_null"`
	Block          uint64
	BlockCreatedAt int64
	Asset          string
	Amount         string
	CreatedAt      int64 `gorm:"auto_create_time"`
}

// func NormalizeTransaction(tx types.Tx, chain types.ChainType) *Transaction {
// 	transaction := Transaction{
// 		Hash:           tx.ID,
// 		Chain:          chain,
// 		Sender:         tx.From,
// 		Recipient:      tx.To,
// 		Block:          tx.Block,
// 		BlockCreatedAt: tx.BlockCreatedAt,
// 	}

// 	switch tx.Type {
// 	case types.TxTransfer:
// 		meta, ok := tx.Metadata.(*types.Transfer)
// 		if !ok {
// 			log.WithFields(log.Fields{
// 				"tx_hash": tx.ID,
// 				"chain":   chain,
// 				"meta":    tx.Metadata,
// 			}).Info("Casting to types.TxTransfer error")

// 			return nil
// 		}

// 		transaction.Asset = meta.Asset
// 		transaction.Amount = string(meta.Value)
// 	default:
// 		return nil
// 	}

// 	return &transaction
// }
