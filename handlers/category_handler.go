package handlers

import (
	"awesomeProject3/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	userID, _ := c.Get("userID")
	var categoryList []models.Category
	DB.Where("user_id = ?", userID).Preload("Transactions").Find(&categoryList)
	c.JSON(http.StatusOK, categoryList)
}

func AddCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists models.Category
	err := DB.Where("name = ? AND user_id = ?", category.Name, userID).First(&exists).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Категория с таким названием уже существует"})
		return
	}

	category.UserID = userID.(uint)
	DB.Create(&category)
	c.JSON(http.StatusCreated, category)
}

func GetCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.Category

	if err := DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("Transactions").
		Preload("Transactions.Tags").
		First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))

	if err := DB.Where("id = ? AND user_id = ?", id, userID).First(&models.Category{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	DB.Where("user_id = ?", userID).Delete(&models.Category{}, id)
	c.Status(http.StatusNoContent)
}
