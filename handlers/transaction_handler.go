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
	var transactionList []models.Transaction

	query := DB.Preload("Category").Preload("Tags")
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
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create transaction"})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}

func GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var transaction models.Transaction
	if err := DB.Preload("Category").Preload("Tags").First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}
func UpdateTransaction(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var transaction models.Transaction
	if err := DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction.ID = uint(id)
	DB.Save(&transaction)
	c.JSON(http.StatusOK, transaction)
}
func DeleteTransaction(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var transaction models.Transaction

	if err := DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	DB.Delete(&transaction)
	c.Status(http.StatusNoContent)
}
