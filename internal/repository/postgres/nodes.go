package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

func (d *Database) InsertNodes(ctx context.Context, nodes []models.Node) error {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&nodes).Error; err != nil {
		return fmt.Errorf("failed to insert nodes: %w", err)
	}

	return nil
}

func (d *Database) GetNodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	if err := d.Gorm.
		WithContext(ctx).
		Order("chain, id asc").
		Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("failed to find nodes: %w", err)
	}

	return nodes, nil
}

func (d *Database) GetNodesByChain(ctx context.Context, chain string) ([]models.Node, error) {
	var nodes []models.Node
	if err := d.Gorm.
		WithContext(ctx).
		Where("chain = ?", chain).
		Order("chain, id asc").
		Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("failed to find nodes by chain: %w", err)
	}

	return nodes, nil
}
