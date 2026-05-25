package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

type LocalTx struct {
	ID         uint    `gorm:"primaryKey"`
	Amount     float64 `json:"amount"`
	CategoryID uint    `json:"category_id"`
	Note       string  `json:"note"`
	UserID     uint    `json:"user_id"`
}

func AddTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Amount     float64 `json:"amount"`
		CategoryID uint    `json:"category_id"`
		Note       string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Table("user_transactions").AutoMigrate(&LocalTx{})

	tx := LocalTx{
		Amount:     input.Amount,
		CategoryID: 1,
		Note:       input.Note,
		UserID:     userID.(uint),
	}

	if err := DB.Table("user_transactions").Create(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func GetTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	var transactionList []LocalTx

	DB.Table("user_transactions").AutoMigrate(&LocalTx{})

	if err := DB.Table("user_transactions").Where("user_id = ?", userID).Order("id desc").Find(&transactionList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactionList)
}

func GetTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var tx LocalTx
	if err := DB.Table("user_transactions").Where("id = ? AND user_id = ?", id, userID).First(&tx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func UpdateTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var tx LocalTx
	if err := DB.Table("user_transactions").Where("id = ? AND user_id = ?", id, userID).First(&tx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx.ID = uint(id)
	tx.UserID = userID.(uint)
	DB.Table("user_transactions").Save(&tx)
	c.JSON(http.StatusOK, tx)
}

func DeleteTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := DB.Table("user_transactions").Where("id = ? AND user_id = ?", id, userID).Delete(&LocalTx{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
