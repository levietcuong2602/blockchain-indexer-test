package models

// Event - Events emitted from smart contracts to be held in this table
type Event struct {
	BlockHash       string   `gorm:"column:blockhash;type:char(66);not null;primaryKey"`
	Index           uint     `gorm:"column:index;type:integer;not null;primaryKey"`
	Origin          string   `gorm:"column:origin;type:char(42);not null;index"`
	Topics          []string `gorm:"column:topics;type:text[];not null;index:,type:gin"`
	Data            []byte   `gorm:"column:data;type:bytea"`
	TransactionHash string   `gorm:"column:txhash;type:char(66);not null;index"`
}
