package handlers

import (
	"net/http"
	"strconv"

	"github.com/Gulner-GI/ToDoApp/models"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, models.Todos)
}

func CreateTodo(c *gin.Context) {
	var newTodo models.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTodo.ID = len(models.Todos) + 1
	models.Todos = append(models.Todos, newTodo)
	c.JSON(http.StatusCreated, newTodo)
}

func UpdateTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTodo models.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, t := range models.Todos {
		if t.ID == id {
			models.Todos[i].Task = updatedTodo.Task
			models.Todos[i].Done = updatedTodo.Done
			c.JSON(http.StatusOK, models.Todos[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

func DeleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, t := range models.Todos {
		if t.ID == id {
			models.Todos = append(models.Todos[:i], models.Todos[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}
