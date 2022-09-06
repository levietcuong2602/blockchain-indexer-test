package repository

import (
	"context"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"
)

type Storage interface {
	// Block trackers
	GetBlockTracker(ctx context.Context, chain types.ChainType) (*models.BlockTracker, error)
	UpsertBlockTracker(ctx context.Context, chain types.ChainType, height int64) error

	// Transactions
	InsertTransactions(ctx context.Context, txs []models.Transaction) error
}
