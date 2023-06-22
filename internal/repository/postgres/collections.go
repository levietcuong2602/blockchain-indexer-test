package postgres

import (
	"context"
	"fmt"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"gorm.io/gorm/clause"
)

func (d *Database) InsertCollections(ctx context.Context, collections []models.Collection) error {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&collections).Error; err != nil {
		return fmt.Errorf("failed to insert collections: %w", err)
	}

	return nil
}

func (d *Database) GetCollections(ctx context.Context) ([]models.Collection, error) {
	var collections []models.Collection
	if err := d.Gorm.
		WithContext(ctx).
		Order("chain, id asc").
		Find(&collections).Error; err != nil {
		return nil, fmt.Errorf("failed to find collections: %w", err)
	}

	return collections, nil
}
