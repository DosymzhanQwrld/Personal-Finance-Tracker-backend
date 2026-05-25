package config

import (
	"awesomeProject3/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := "host=db user=postgres password=postgres dbname=finance_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&models.User{}, &models.Category{}, &models.Tag{}, &models.Transaction{})
	return db
}
