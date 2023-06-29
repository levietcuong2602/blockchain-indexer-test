package postgres

import (
	"context"
	"fmt"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"gorm.io/gorm/clause"
)

func (d *Database) InsertNftBalance(ctx context.Context, nftBalance models.NftBalance) (*models.NftBalance, error) {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&nftBalance).Error; err != nil {
		return nil, fmt.Errorf("failed to insert nft_balance: %w", err)
	}

	return &nftBalance, nil
}

func (d *Database) FindNftBalanceByOwnerContractAndTokenId(ctx context.Context, owner, contract string, tokenId uint64) (*models.NftBalance, error) {
	var nftBalance models.NftBalance

	if err := d.Gorm.
		WithContext(ctx).
		Where("owner = ?", owner).
		Where("contract = ?", contract).
		Where("token_id = ?", tokenId).
		First(&nftBalance).Error; err != nil {
		return nil, fmt.Errorf("failed to find nft_balance: %w", err)
	}

	return &nftBalance, nil
}

func (d *Database) UpdateNftBalance(ctx context.Context, nftBalanceUpdate models.NftBalance) (*models.NftBalance, error) {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Updates(&nftBalanceUpdate).Error; err != nil {
		return nil, fmt.Errorf("failed to insert nft_balance: %w", err)
	}

	return &nftBalanceUpdate, nil
}
