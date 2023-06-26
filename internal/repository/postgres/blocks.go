package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

func (d *Database) InsertBlock(ctx context.Context, block models.Block) error {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&block).Error; err != nil {
		return fmt.Errorf("failed to insert blocks: %w", err)
	}

	return nil
}

func (d *Database) GetBlocks(ctx context.Context) ([]models.Block, error) {
	var blocks []models.Block
	if err := d.Gorm.
		WithContext(ctx).
		Order("chain, id asc").
		Find(&blocks).Error; err != nil {
		return nil, fmt.Errorf("failed to find blocks: %w", err)
	}

	return blocks, nil
}

func (d *Database) GetBlocksByChain(ctx context.Context, chain string) ([]models.Block, error) {
	var blocks []models.Block
	if err := d.Gorm.
		WithContext(ctx).
		Where("chain = ?", chain).
		Order("chain, id asc").
		Find(&blocks).Error; err != nil {
		return nil, fmt.Errorf("failed to find blocks by chain: %w", err)
	}

	return blocks, nil
}
