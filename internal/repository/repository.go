package repository

import (
	"context"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

type Storage interface {
	// Block trackers
	GetBlockTracker(ctx context.Context, chain string) (*models.BlockTracker, error)
	UpsertBlockTracker(ctx context.Context, chain string, height int64) error

	// Transactions
	InsertTransactions(ctx context.Context, txs []models.Transaction) error
}
