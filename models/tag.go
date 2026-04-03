package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name         string        `json:"name"`
	Transactions []Transaction `gorm:"many2many:transaction_tags;"`
}
