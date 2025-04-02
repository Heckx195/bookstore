package main

import (
	"restapi/config"
	"restapi/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DBInit()

	r := gin.Default()

	// Book Routes
	r.POST("/books", handlers.CreateBook)
	r.GET("/books", handlers.GetBooks)
	r.GET("/books/:id", handlers.GetBooksById)
	r.PUT("/books/:id", handlers.UpdateBook)
	r.DELETE("/books/:id", handlers.DeleteBook)

	// Author Routes
	r.POST("/authors", handlers.CreateAuthor)
	r.GET("/authors", handlers.GetAuthors)

	// Category Routes
	r.POST("/categories", handlers.CreateCategory)
	r.GET("/categories", handlers.GetCategories)

	r.Run(":8080")
}
