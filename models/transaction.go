package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount     float64  `json:"amount"`
	Note       string   `json:"note"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"-"`
	Tags       []Tag    `gorm:"many2many:transaction_tags;"`
}
