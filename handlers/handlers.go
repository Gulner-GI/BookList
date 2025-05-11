package handlers

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/Gulner-GI/ToDoApp/models"

	"github.com/gin-gonic/gin"
)

func FindToDo(id int) (idx int) {
	for i := range models.Todos {
		if models.Todos[i].ID == id {
			return i
		}
	}
	return -1
}

func GetTodos(c *gin.Context) {
	idParam := c.Query("id")
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		idx := FindToDo(id)
		if idx != -1 {
			c.JSON(http.StatusOK, models.Todos[idx])
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
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
	var updatedTodo models.Update
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idx := FindToDo(id)
	if idx != -1 {
		if updatedTodo.Task != nil {
			models.Todos[idx].Task = *updatedTodo.Task
		}
		if updatedTodo.Done != nil {
			models.Todos[idx].Done = *updatedTodo.Done
		}
		c.JSON(http.StatusOK, models.Todos[idx])
		return
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
	idx := FindToDo(id)
	if idx != -1 {
		models.Todos = slices.Delete(models.Todos, idx, idx+1)
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}
