package postgres

import (
	"context"
	"fmt"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"gorm.io/gorm/clause"
)

func (d *Database) InsertCollection(ctx context.Context, collection models.Collection) (*models.Collection, error) {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&collection).Error; err != nil {
		return nil, fmt.Errorf("failed to insert collections: %w", err)
	}

	return &collection, nil
}

func (d *Database) GetCollections(ctx context.Context, name string, page, limit int, recent bool) ([]models.Collection, error) {
	collections := make([]models.Collection, 0)

	orderBy := "asc"
	if recent {
		orderBy = "desc"
	}

	if err := d.Gorm.
		WithContext(ctx).
		Where("name like ?", name).
		Order("created_at " + orderBy).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&collections).Error; err != nil {
		return nil, fmt.Errorf("failed to find collections: %w", err)
	}

	return collections, nil
}

func (d *Database) GetCollectionTotalCount(ctx context.Context, name string) (int64, error) {
	var totalCount int64

	if err := d.Gorm.
		WithContext(ctx).
		Model(&models.Collection{}).
		Where("name = ?", name).
		Count(&totalCount).Error; err != nil {
		return totalCount, fmt.Errorf("failed to get count of collection: %w", err)
	}

	return totalCount, nil
}

func (d *Database) FindCollectionByContract(ctx context.Context, contract string) (*models.Collection, error) {
	var collection models.Collection

	if err := d.Gorm.
		WithContext(ctx).
		Where("contract = ?", contract).
		First(&collection).Error; err != nil {
		return nil, fmt.Errorf("failed to find collection: %w", err)
	}

	return &collection, nil
}
