package main

import (
	"awesomeProject3/config"
	"awesomeProject3/handlers"
	"awesomeProject3/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()
	handlers.DB = db
	r := gin.Default()
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", handlers.GetUsers)
		protected.GET("/transactions", handlers.GetTransactions)
		protected.POST("/transactions", handlers.AddTransaction)
		protected.GET("/transactions/:id", handlers.GetTransaction)
		protected.PUT("/transactions/:id", handlers.UpdateTransaction)
		protected.DELETE("/transactions/:id", handlers.DeleteTransaction)
		protected.GET("/categories", handlers.GetCategories)
		protected.POST("/categories", handlers.AddCategory)
		protected.GET("/categories/:id", handlers.GetCategory)
		protected.DELETE("/categories/:id", handlers.DeleteCategory)
		protected.GET("/tags", handlers.GetTags)
		protected.POST("/tags", handlers.AddTag)
		protected.POST("/transactions/:id/tags", handlers.AttachTag)
	}
	r.Run(":8080")
}
