package handlers

import (
	"net/http"
	"restapi/models"

	"github.com/gin-gonic/gin"
)

var categories = []models.Category{
	{ID: 1, Name: "Fiction"},
	{ID: 2, Name: "Fantasy"},
	{ID: 3, Name: "Drama"},
	{ID: 4, Name: "Romantic"},
}

// POST: /categories
func CreateCategory(c *gin.Context) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)
	c.JSON(http.StatusCreated, newCategory)
}

// GET: /categories
func GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, categories)
}

// Getter categories
func getCategories() []models.Category {
	return categories
}
