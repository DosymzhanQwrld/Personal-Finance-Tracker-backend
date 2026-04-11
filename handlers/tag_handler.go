package handlers

import (
	"awesomeProject3/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	userID, _ := c.Get("userID")
	var tagList []models.Tag
	DB.Where("user_id = ?", userID).Preload("Transactions").Find(&tagList)
	c.JSON(http.StatusOK, tagList)
}

func AddTag(c *gin.Context) {
	userID, _ := c.Get("userID")
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists models.Tag
	err := DB.Where("name = ? AND user_id = ?", tag.Name, userID).First(&exists).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Тег с таким названием уже существует"})
		return
	}

	tag.UserID = userID.(uint)
	DB.Create(&tag)
	c.JSON(http.StatusCreated, tag)
}

func AttachTag(c *gin.Context) {
	userID, _ := c.Get("userID")
	txID, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		TagID uint `json:"tag_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var transaction models.Transaction
	var tag models.Tag

	if err := DB.Where("id = ? AND user_id = ?", txID, userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found or access denied"})
		return
	}

	if err := DB.Where("id = ? AND user_id = ?", input.TagID, userID).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found or access denied"})
		return
	}

	DB.Model(&transaction).Association("Tags").Append(&tag)
	c.JSON(http.StatusOK, gin.H{"message": "Tag attached successfully"})
}
