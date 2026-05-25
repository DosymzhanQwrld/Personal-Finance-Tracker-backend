package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount     float64  `gorm:"not null" json:"amount"`
	Note       string   `json:"note"`
	UserID     uint     `gorm:"not null" json:"user_id"`
	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
	Tags       []Tag    `gorm:"many2many:transaction_tags;" json:"tags"`
}
