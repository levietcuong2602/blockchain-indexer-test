package models

// Block - Mined block info holder table model
type Block struct {
	Hash                string      `gorm:"column:hash;type:char(66);primaryKey"`
	Number              uint64      `gorm:"column:number;type:bigint;not null;unique;index:,sort:asc"`
	Time                uint64      `gorm:"column:time;type:bigint;not null;index:,sort:asc"`
	ParentHash          string      `gorm:"column:parent_hash;type:char(66);not null"`
	Difficulty          string      `gorm:"column:difficulty;type:varchar;not null"`
	GasUsed             uint64      `gorm:"column:gas_used;type:bigint;not null"`
	GasLimit            uint64      `gorm:"column:gas_limit;type:bigint;not null"`
	Nonce               string      `gorm:"column:nonce;type:varchar;not null"`
	Miner               string      `gorm:"column:miner;type:char(42);not null"`
	Size                float64     `gorm:"column:size;type:float(8);not null"`
	StateRootHash       string      `gorm:"column:state_root_hash;type:char(66);not null"`
	UncleHash           string      `gorm:"column:uncle_hash;type:char(66);not null"`
	TransactionRootHash string      `gorm:"column:tx_root_hash;type:char(66);not null"`
	ReceiptRootHash     string      `gorm:"column:receipt_root_hash;type:char(66);not null"`
	ExtraData           []byte      `gorm:"column:extra_data;type:bytea"`
	Transactions        Transaction `gorm:"foreignKey:blockhash"`
	Events              Event       `gorm:"foreignKey:blockhash"`
}
