package routes

import (
	"github.com/Gulner-GI/BookList/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/books", handlers.FindBooks)
	r.HEAD("/books", handlers.HeadBooks)
	r.OPTIONS("/books", handlers.Options)
	r.POST("/books", handlers.AddBook)
	r.PATCH("/books/:id", handlers.UpdateBook)
	r.DELETE("/books/:id", handlers.DeleteBook)
	return r
}
