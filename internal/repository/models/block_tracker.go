package models

import "github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"

type BlockTracker struct {
	Chain     types.ChainType `gorm:"primary_key:true; type:varchar(64)"`
	Height    int64           `gorm:"default:0"`
	CreatedAt int64           `gorm:"auto_create_time"`
	UpdatedAt int64           `gorm:"auto_update_time"`
}
