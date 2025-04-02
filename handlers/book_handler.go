package handlers

import (
	"fmt"
	"net/http"
	"restapi/config"
	"restapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST: /books
func CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valide, msg := ValidateInput(newBook)
	if !valide {
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	if err := config.DB.Create(&newBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

// GET: /books
func GetBooks(c *gin.Context) {
	// Get all books.
	var books []models.Book
	if err := config.DB.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	// Filter & Pagination logic.
	category := c.Query("category")
	pageParam := c.Query("page")
	limitParam := c.Query("limit")

	page, err_page := strconv.Atoi(pageParam)
	limit, err_limit := strconv.Atoi(limitParam)

	var checkCategory bool = true
	if category == "" {
		fmt.Println("Info: ategory is empty")
		checkCategory = false
	}

	// Set defaults
	if err_page != nil || page < 1 {
		fmt.Println("Page param is nil or smaller 1 // Set default value 1")
		page = 1
	}
	if err_limit != nil {
		fmt.Println("Limit param is nil // Set default value 2")
		limit = len(books)
	}

	// Get category name by id
	var catId int = -1
	if checkCategory {
		categories, err := GetAllCategories()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

		for _, value := range categories {
			if value.Name == category {
				catId = value.ID
				break
			}
		}
	}

	// Filter books by category.
	var selectedBooks []models.Book
	for _, book := range books {
		fmt.Println("Book:", book)
		if checkCategory && book.CategoryID == catId {
			selectedBooks = append(selectedBooks, book)
		} else if !checkCategory {
			selectedBooks = append(selectedBooks, book)
		}
	}

	// Paginate the selected books.
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	// Ensure indices are within bounds
	if startIndex >= len(selectedBooks) {
		c.JSON(http.StatusOK, []models.Book{}) // Return empty if page exceeds available books
		return
	}
	if endIndex > len(selectedBooks) {
		endIndex = len(selectedBooks)
	}

	pagedBooks := selectedBooks[startIndex:endIndex]
	c.JSON(http.StatusOK, pagedBooks)
}

// GET: /books/:id
func GetBooksById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// PUT: /books/:id
func UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.Book{}).Where("id = ?", id).Updates(updatedBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

// DELETE: /books/:id
func DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := config.DB.Delete(&models.Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

// Private helper functions.
func ValidateInput(b models.Book) (bool, string) {
	// Check minimum price
	if b.Price < 1_000 {
		return false, "Price too low"
	}
	// Check all fields filled
	if b.Title == "" || b.AuthorID == 0 || b.CategoryID == 0 || b.Price == 0 {
		return false, "All fields must be filled"
	}

	return true, "valide"
}

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
