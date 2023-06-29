package postgres

import (
	"context"
	"fmt"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"gorm.io/gorm/clause"
)

func (d *Database) InsertEvents(ctx context.Context, txs []models.Event) error {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&txs, batchSize).Error; err != nil {
		return fmt.Errorf("failed to insert events: %w", err)
	}

	return nil
}
