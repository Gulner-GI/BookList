package routes

import (
	_ "github.com/Gulner-GI/BookList/docs"
	"github.com/Gulner-GI/BookList/handlers"
	"github.com/Gulner-GI/BookList/loggers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(loggers.LogFile))
	r.GET("/books", handlers.FindBooks)
	r.HEAD("/books", handlers.HeadBooks)
	r.OPTIONS("/books", handlers.Options)
	r.POST("/books", handlers.AddBook)
	r.PATCH("/books/:id", handlers.UpdateBook)
	r.DELETE("/books/:id", handlers.DeleteBook)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
