package main

import (
	"awesomeProject3/config"
	"awesomeProject3/handlers"
	"awesomeProject3/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()

	err := db.AutoMigrate(&models.Category{}, &models.Tag{}, &models.Transaction{})
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

	r.GET("/tags", handlers.GetTags)
	r.POST("/tags", handlers.AddTag)
	r.POST("/transactions/:id/tags", handlers.AttachTag)

	r.Run(":8080")
}
