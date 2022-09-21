package models

type Node struct {
	ID         int64 `gorm:"primary_key; auto_increment"`
	Chain      string
	Scheme     string
	Host       string `gorm:"type:varchar(256);uniqueIndex"`
	Enabled    bool   `gorm:"default:true"`
	Monitoring bool   `gorm:"default:true"`
}
