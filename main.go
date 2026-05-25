package main

import (
	"awesomeProject3/config"
	"awesomeProject3/handlers"
	"awesomeProject3/middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	db := config.ConnectDatabase()
	handlers.DB = db

	r := gin.Default()

	r.Use(CORSMiddleware())

	r.StaticFile("/", "./frontend/index.html")

	client := resty.New()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.GET("/exchange-rate", func(c *gin.Context) {
		exchangeURL := os.Getenv("EXCHANGE_SERVICE_URL")
		if exchangeURL == "" {
			exchangeURL = "http://localhost:8081"
		}

		resp, err := client.R().
			SetQueryParams(map[string]string{"currency": "USD"}).
			Get(exchangeURL + "/rate")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Exchange service unreachable",
				"details": err.Error(),
			})
			return
		}
		c.Data(http.StatusOK, "application/json", resp.Body())
	})

	r.POST("/send-notify", func(c *gin.Context) {
		notifyURL := os.Getenv("NOTIFY_SERVICE_URL")
		if notifyURL == "" {
			notifyURL = "http://localhost:8082"
		}

		var msg map[string]interface{}
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		resp, err := client.R().
			SetBody(msg).
			Post(notifyURL + "/send")

		if err != nil {
			c.Status(http.StatusBadGateway)
			return
		}

		c.Data(resp.StatusCode(), "application/json", resp.Body())
	})

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
	fmt.Println("Finance App is running on :8080")
	r.Run(":8080")
}
