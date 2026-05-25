package handlers

import (
	"awesomeProject3/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func AddTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Amount     float64 `json:"amount" binding:"required"`
		CategoryID uint    `json:"category_id" binding:"required"`
		Note       string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := models.Transaction{
		Amount:     input.Amount,
		CategoryID: input.CategoryID,
		Note:       input.Note,
		UserID:     userID.(uint),
	}

	if err := DB.Create(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func GetTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	var transactionList []models.Transaction

	if err := DB.Where("user_id = ?", userID).Preload("Tags").Preload("Category").Order("id desc").Find(&transactionList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactionList)
}

func GetTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	var tx models.Transaction
	if err := DB.Where("id = ? AND user_id = ?", id, userID).Preload("Tags").First(&tx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func UpdateTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	var tx models.Transaction
	if err := DB.Where("id = ? AND user_id = ?", id, userID).First(&tx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.ID = uint(id)
	tx.UserID = userID.(uint)

	DB.Save(&tx)
	c.JSON(http.StatusOK, tx)
}

func DeleteTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
