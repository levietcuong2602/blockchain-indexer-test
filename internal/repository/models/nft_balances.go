package models

type NftBalance struct {
	ID       int64  `gorm:"primary_key; auto_increment"`
	Contract string `gorm:"type:varchar(256); not_null"`
	TokenId  uint64 `gorm:"type:bigint;not null;index"`
	Owner    string `gorm:"type:varchar(256); not_null"`
	Amount   int64  `gorm:"default:0"`
}
