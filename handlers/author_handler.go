package handlers

import (
	"net/http"
	"restapi/models"

	"github.com/gin-gonic/gin"
)

var authors = []models.Author{
	{ID: 1, Name: "Athur"},
	{ID: 2, Name: "Toni"},
	{ID: 3, Name: "Peter"},
	{ID: 4, Name: "Max"},
}

// POST: /authors
func CreateAuthor(c *gin.Context) {
	var newAuthor models.Author
	if err := c.ShouldBindJSON(&newAuthor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAuthor.ID = len(authors) + 1
	authors = append(authors, newAuthor)
	c.JSON(http.StatusCreated, newAuthor)
}

// GET: /authors
func GetAuthors(c *gin.Context) {
	c.JSON(http.StatusOK, authors)
}
