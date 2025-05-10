package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	var todos = []Todo{
		{ID: 1, Task: "Learn Go"},
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, todos)
	})
	r.POST("/todos", func(c *gin.Context) {
		var newTodo Todo
		if err := c.BindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo.ID = len(todos) + 1
		todos = append(todos, newTodo)
		c.JSON(http.StatusCreated, newTodo)
	})
	r.Run()
}
