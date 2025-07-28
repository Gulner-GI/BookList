package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Gulner-GI/ToDoApp/db"
	"github.com/Gulner-GI/ToDoApp/models"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	idParam := c.Query("id")
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		row := db.DB.QueryRow("SELECT id, task, done FROM todos WHERE id = ?", id)
		var todo models.Todo
		err = row.Scan(&todo.ID, todo.Task)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todo)
		return
	}
	rows, err := db.DB.Query("SELECT id, task, done FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Task, &todo.Done)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, todos)
}

func HeadTodos(c *gin.Context) {
	GetTodos(c)
	c.Writer.WriteHeaderNow()
	c.Writer.Flush()
}

func OptionsTodos(c *gin.Context) {
	c.Header("Allow", "GET, HEAD, POST, PATCH, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PATCH, DELETE, OPTIONS")
	c.Status(http.StatusOK)
}

func CreateTodo(c *gin.Context) {
	var newTodo models.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.DB.Exec("INSERT INTO todos (task, done) VALUES (?, ?)", newTodo.Task, newTodo.Done)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert todo"})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inserted ID"})
		return
	}
	newTodo.ID = int(id)
	c.JSON(http.StatusCreated, newTodo)
}

func UpdateTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var update models.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if update.Task == nil && update.Done == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}
	row := db.DB.QueryRow("SELECT id, task, done FROM todos WHERE id = ?", id)
	var existingID int
	var currentTask string
	var currentDone bool
	if err := row.Scan(&existingID, &currentTask, &currentDone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	if update.Task != nil {
		currentTask = *update.Task
	}
	if update.Done != nil {
		currentDone = *update.Done
	}
	_, err = db.DB.Exec(
		"UPDATE todos SET task = ?, done = ? WHERE id = ?",
		currentTask, currentDone, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"task": currentTask,
		"done": currentDone,
	})
}

func DeleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check deletion result"})
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
