package models

import (
	"database/sql/driver"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Network string

const (
	MUMBAI   Network = "mumbai"
	ETHEREUM Network = "ethereum"
	COSMOS   Network = "cosmos"
	BINANCE  Network = "binance"
	SOLANA   Network = "solana"
	NEAR     Network = "near"
)

func (n *Network) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		*n = Network(b)
	}
	return nil
}

func (n Network) Value() (driver.Value, error) {
	return string(n), nil
}

type ContractStandard string

const (
	ERC20   Network = "ERC20"
	ERC721  Network = "ERC721"
	ERC1155 Network = "ERC1155"
	UNKNOWN Network = "UNKNOWN"
)

type Collection struct {
	ID              int64 `gorm:"primary_key; auto_increment"`
	Slug            string
	Chain           Network          `sql:"type:network; not_null" gorm:"column:chain"`
	Name            string           `gorm:"type:varchar(256); not_null"`
	Standard        ContractStandard `sql:"type:contract_standard" gorm:"column:standard"`
	Metadata        postgres.Jsonb
	Contract        string `gorm:"type:varchar(256); not_null"`
	TokenCount      int64  `gorm:"default:0"`
	MintedTimestamp int64  `gorm:"auto_create_time"`
	CreatedAt       int64  `gorm:"auto_create_time"`
	UpdatedAt       int64  `gorm:"auto_update_time"`
}
