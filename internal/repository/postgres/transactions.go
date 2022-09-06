package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

const batchSize = 1000

func (d *Database) InsertTransactions(ctx context.Context, txs []models.Transaction) error {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&txs, batchSize).Error; err != nil {
		return fmt.Errorf("failed to insert txs: %w", err)
	}

	return nil
}
