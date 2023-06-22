package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Collection struct {
	ID              int64 `gorm:"primary_key; auto_increment"`
	Slug            string
	Name            string
	Metadata        postgres.Jsonb
	Contract        string `gorm:"type:varchar(256); not_null"`
	TokenCount      int64  `gorm:"default:0"`
	MintedTimestamp int64  `gorm:"auto_create_time"`
	CreatedAt       int64  `gorm:"auto_create_time"`
	UpdatedAt       int64  `gorm:"auto_update_time"`
}
