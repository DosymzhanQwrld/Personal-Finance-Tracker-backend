package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/rate", func(c *gin.Context) {
		currency := c.DefaultQuery("currency", "USD")
		c.JSON(http.StatusOK, gin.H{
			"currency": currency,
			"rate":     448.5,
			"source":   "Exchange Service v1",
		})
	})

	r.Run(":8081")
}
