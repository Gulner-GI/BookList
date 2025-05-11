package routes

import (
	"github.com/Gulner-GI/ToDoApp/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/Todos", handlers.GetTodos)
	r.POST("/Todos", handlers.CreateTodo)
	r.PATCH("/Todos/:id", handlers.UpdateTodo)
	r.DELETE("/Todos/:id", handlers.DeleteTodo)
	return r
}
