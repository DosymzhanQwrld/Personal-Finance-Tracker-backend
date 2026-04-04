package handlers

import (
	"awesomeProject3/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	var tagList []models.Tag
	DB.Preload("Transactions").Find(&tagList)
	c.JSON(http.StatusOK, tagList)
}

func AddTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&tag)
	c.JSON(http.StatusCreated, tag)
}

func AttachTag(c *gin.Context) {
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

	if err := DB.First(&transaction, txID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	if err := DB.First(&tag, input.TagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}
	DB.Model(&transaction).Association("Tags").Append(&tag)

	c.JSON(http.StatusOK, gin.H{"message": "Tag attached successfully"})
}
