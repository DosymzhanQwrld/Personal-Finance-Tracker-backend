package main

import (
	"awesomeProject3/handlers"
	"awesomeProject3/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=finance_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&models.Category{}, &models.Tag{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	handlers.DB = db

	r := gin.Default()

	r.GET("/transactions", handlers.GetTransactions)
	r.POST("/transactions", handlers.AddTransaction)
	r.GET("/transactions/:id", handlers.GetTransaction)
	r.PUT("/transactions/:id", handlers.UpdateTransaction)
	r.DELETE("/transactions/:id", handlers.DeleteTransaction)

	r.GET("/categories", handlers.GetCategories)
	r.POST("/categories", handlers.AddCategory)
	r.GET("/categories/:id", handlers.GetCategory)
	r.DELETE("/categories/:id", handlers.DeleteCategory)

	r.GET("/tags", handlers.GetTags)                     // 10. Все теги
	r.POST("/tags", handlers.AddTag)                     // 11. Создать тег
	r.POST("/transactions/:id/tags", handlers.AttachTag) // 12. Привязать тег к транзакции

	r.Run(":8080")
}
