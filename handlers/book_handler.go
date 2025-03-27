package handlers

import (
	"fmt"
	"net/http"
	"restapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var books = []models.Book{
	{ID: 0, Title: "Book One", AuthorID: 1, CategoryID: 1, Price: 2000.00},
	{ID: 1, Title: "Book Two", AuthorID: 2, CategoryID: 2, Price: 1900.00},
	{ID: 2, Title: "Book Three", AuthorID: 1, CategoryID: 3, Price: 2998.99},
	{ID: 3, Title: "Book Four", AuthorID: 3, CategoryID: 1, Price: 1200.00},
	{ID: 4, Title: "Book Five", AuthorID: 2, CategoryID: 2, Price: 2400.00},
}

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

	newBook.ID = len(books) //  + 1
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

// GET: /books
func GetBooks(c *gin.Context) {
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
		categories := getCategories()
		for _, value := range categories {
			if value.Name == category {
				catId = value.ID
				break
			}
		}
	}

	// Filter
	var selectedBooks []models.Book
	for _, book := range books {
		fmt.Println("Book:", book)
		if checkCategory && book.CategoryID == catId {
			selectedBooks = append(selectedBooks, book)
		} else if !checkCategory {
			selectedBooks = append(selectedBooks, book)
		}
	}

	// Paginate the selected books
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

	for _, book := range books {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
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

	for i, book := range books {
		if book.ID == id {
			books[i] = updatedBook
			books[i].ID = id
			c.JSON(http.StatusOK, books[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

// DELETE: /books/:id
func DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	for i, book := range books {
		if book.ID == id {
			// Copy all books from start till i and from i+1 till end.
			books = append(books[:i], books[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
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
