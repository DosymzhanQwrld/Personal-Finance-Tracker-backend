package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name         string        `gorm:"uniqueIndex:idx_user_tag" json:"name"`
	UserID       uint          `gorm:"uniqueIndex:idx_user_tag" json:"user_id"`
	Transactions []Transaction `gorm:"many2many:transaction_tags;"`
}
