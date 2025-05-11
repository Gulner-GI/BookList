package main

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

func findtodo(id int) int {
	for i, t := range todos {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func removetodo(idx int) {
	todos = slices.Delete(todos, idx, idx+1)
}

func main() {
	r := gin.Default()
	r.GET("/todos/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		idx := findtodo(id)
		if idx == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusOK, todos[idx])
	})

	r.POST("/todos", func(c *gin.Context) {
		var newTodo Todo
		if err := c.BindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo.ID = len(todos) + 1
		newTodo.Done = false
		todos = append(todos, newTodo)
		c.JSON(http.StatusCreated, newTodo)
	})

	r.PATCH("/todos/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var updatedTodo struct {
			Task *string `json:"task,omitempty"`
			Done *bool   `json:"done,omitempty"`
		}
		if err := c.BindJSON(&updatedTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		idx := findtodo(id)
		if idx == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "todo not found"})
			return
		}
		if updatedTodo.Task != nil {
			todos[idx].Task = *updatedTodo.Task
		}
		if updatedTodo.Done != nil {
			todos[idx].Done = *updatedTodo.Done
		}
		c.JSON(http.StatusOK, todos[idx])
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ivalid id"})
			return
		}
		idx := findtodo(id)
		if idx == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "todo not found"})
			return
		}
		removetodo(idx)
		c.Status(http.StatusNoContent)
	})

	r.Run()
}
