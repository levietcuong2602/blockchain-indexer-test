package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

func (d *Database) GetBlockTracker(ctx context.Context, chain string) (*models.BlockTracker, error) {
	var tracker models.BlockTracker

	if err := d.Gorm.
		WithContext(ctx).
		Where("chain = ?", chain).
		First(&tracker).Error; err != nil {
		return nil, fmt.Errorf("failed to find a block tracker: %w", err)
	}

	return &tracker, nil
}

func (d *Database) UpsertBlockTracker(ctx context.Context, chain string, height int64) error {
	tracker := models.BlockTracker{Chain: chain, Height: height}

	if err := d.Gorm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain"}},
		DoUpdates: clause.AssignmentColumns([]string{"height", "updated_at"}),
	}).Create(&tracker).Error; err != nil {
		return fmt.Errorf("failed to upsert block tracker: %w", err)
	}

	return nil
}
