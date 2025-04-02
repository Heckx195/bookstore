package handlers

import (
	"net/http"
	"restapi/config"
	"restapi/models"

	"github.com/gin-gonic/gin"
)

// POST: /authors
func CreateAuthor(c *gin.Context) {
	var newAuthor models.Author
	if err := c.ShouldBindJSON(&newAuthor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&newAuthor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.JSON(http.StatusCreated, newAuthor)
}

// GET: /authors
func GetAuthors(c *gin.Context) {
	var authors []models.Author
	if err := config.DB.Find(&authors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch authors"})
		return
	}

	c.JSON(http.StatusOK, authors)
}
