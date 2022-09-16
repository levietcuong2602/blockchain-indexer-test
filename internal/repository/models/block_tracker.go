package models

type BlockTracker struct {
	Chain     string `gorm:"primary_key:true; type:varchar(64)"`
	Height    int64  `gorm:"default:0"`
	CreatedAt int64  `gorm:"auto_create_time"`
	UpdatedAt int64  `gorm:"auto_update_time"`
}
