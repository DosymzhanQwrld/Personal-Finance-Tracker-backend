package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount     float64   `json:"amount"`
	Note       string    `json:"note"`
	UserID     uint      `json:"user_id"`
	CategoryID uint      `json:"category_id"`
	Category   *Category `json:"category,omitempty"`
	Tags       []Tag     `gorm:"many2many:transaction_tags;" json:"tags,omitempty"`
}
