package models

type TokenBalances struct {
	ID         int64  `gorm:"primary_key; auto_increment"`
	Contract   string `gorm:"type:varchar(256); not_null"`
	TokenId    string `gorm:"type:varchar(256); not_null"`
	Owner      string `gorm:"type:varchar(256); not_null"`
	Amount     int64  `gorm:"default:0"`
	AcquiredAt int64
}
