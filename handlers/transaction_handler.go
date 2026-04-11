package handlers

import (
	"awesomeProject3/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	var transactionList []models.Transaction
	query := DB.Where("user_id = ?", userID).Preload("Category").Preload("Tags")
	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	query = query.Order(sort + " " + order)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * pageSize
	query.Limit(pageSize).Offset(offset).Find(&transactionList)
	c.JSON(http.StatusOK, transactionList)
}

func AddTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := DB.Where("id = ? AND user_id = ?", transaction.CategoryID, userID).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Категория не найдена"})
		return
	}

	transaction.UserID = userID.(uint)

	if err := DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create transaction"})
		return
	}

	DB.Preload("Category").Preload("Tags").First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, transaction)
}

func GetTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var transaction models.Transaction
	if err := DB.Where("user_id = ?", userID).Preload("Category").Preload("Tags").First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

func UpdateTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var transaction models.Transaction
	if err := DB.Where("user_id = ?", userID).First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction.ID = uint(id)
	transaction.UserID = userID.(uint)
	DB.Save(&transaction)
	c.JSON(http.StatusOK, transaction)
}

func DeleteTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := DB.Where("user_id = ?", userID).Delete(&models.Transaction{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
