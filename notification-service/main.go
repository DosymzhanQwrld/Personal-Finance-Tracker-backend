package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationRequest struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func main() {
	r := gin.Default()

	r.POST("/send", func(c *gin.Context) {
		var req NotificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("НОВОЕ УВЕДОМЛЕНИЕ: Пользователю %d отправлено сообщение (%s): %s\n",
			req.UserID, req.Type, req.Message)

		c.JSON(http.StatusOK, gin.H{
			"status":  "sent",
			"message": "Notification processed successfully",
		})
	})

	fmt.Println("Notification Service is running on :8082")
	r.Run(":8082")
}
