package handlers

import (
	"fmt"
	"net/http"
	"restapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var books = []models.Book{
	{ID: 1, Title: "Book One", AuthorID: 1, CategoryID: 1, Price: 2000.00},
	{ID: 2, Title: "Book Two", AuthorID: 2, CategoryID: 2, Price: 1900.00},
	{ID: 3, Title: "Book Three", AuthorID: 1, CategoryID: 3, Price: 2998.99},
	{ID: 4, Title: "Book Four", AuthorID: 3, CategoryID: 1, Price: 1200.00},
	{ID: 5, Title: "Book Five", AuthorID: 2, CategoryID: 2, Price: 2400.00},
}

// POST: /books
func CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newBook.ID = len(books) + 1
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

// GET: /books
func GetBooks(c *gin.Context) {
	category := c.Query("category")
	pageParam := c.Query("page")

	page, err_page := strconv.Atoi(pageParam)

	if category == "" {
		fmt.Println("Error category is empty")
	}
	if err_page != nil {
		fmt.Println("Error err_page nil")
	}
	if page > len(books) {
		fmt.Println("Error page too big")
	}

	if category == "" || err_page != nil || page > len(books) {
		fmt.Println("Error in getbooks params")
		c.JSON(http.StatusOK, books)
	}

	// Set default
	if err_page != nil {
		page = 1
	}

	// Pagination and filter
	var selectedBooks []models.Book
	for i := page; i < len(books); i++ {
		book := books[i]
		if strconv.Itoa(book.CategoryID) == category {
			selectedBooks = append(selectedBooks, book)
		}
	}
	c.JSON(http.StatusOK, selectedBooks)
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
	if b.Price > 1_000 {
		return false, "Price too low"
	}
	// Check all fields filled
	if b.Title == "" || b.AuthorID == 0 || b.CategoryID == 0 || b.Price == 0 {
		return false, "All fields must be filled"
	}

	return true, "valide"
}
