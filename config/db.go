package config

import (
	"fmt"
	"log"

	"awesomeProject3/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	dsn := "host=db user=postgres password=postgres dbname=finance_db port=5432 sslmode=disable"
	dbURL := "postgres://postgres:postgres@db:5432/finance_db?sslmode=disable"

	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		dsn = "host=localhost user=postgres password=postgres dbname=finance_db port=5432 sslmode=disable"
		dbURL = "postgres://postgres:postgres@localhost:5432/finance_db?sslmode=disable"

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database: ", err)
		}
	}

	DB = db

	err = db.AutoMigrate(&models.User{}, &models.Category{}, &models.Tag{}, &models.Transaction{})
	if err != nil {
		log.Printf("GORM AutoMigrate warning: %v", err)
	}

	m, err := migrate.New("file://config/migrations", dbURL)
	if err != nil {
		log.Printf("Could not start golang-migrate instance: %v (проверь, есть ли папка config/migrations)", err)
		return db
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Golang-migrate failed: %v", err)
	} else {
		fmt.Println("Database migrated successfully using golang-migrate!")
	}

	return db
}
