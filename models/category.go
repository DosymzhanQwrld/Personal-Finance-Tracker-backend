package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name         string        `gorm:"uniqueIndex:idx_user_category" json:"name"`
	Type         string        `json:"type"`
	UserID       uint          `gorm:"uniqueIndex:idx_user_category" json:"user_id"`
	Transactions []Transaction `json:"transactions"`
}
