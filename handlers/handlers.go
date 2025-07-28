package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Gulner-GI/BookList/db"
	"github.com/Gulner-GI/BookList/models"

	"github.com/gin-gonic/gin"
)

func FindBooks(c *gin.Context) {
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
		row := db.DB.QueryRow("SELECT id, title, year, genre, status, link FROM books WHERE id = ?", id)
		var book models.Book
		err = row.Scan(&book.ID, &book.Title, &book.Year, &book.Genre, &book.Status, &book.Link)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if book.Status {
			book.StatusText = "completed"
		} else {
			book.StatusText = "in process"
		}
		c.IndentedJSON(http.StatusOK, book)
		return
	}
	rows, err := db.DB.Query("SELECT id, title, year, genre, status, link FROM books")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Year, &book.Genre, &book.Status, &book.Link)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if book.Status {
			book.StatusText = "completed"
		} else {
			book.StatusText = "in process"
		}
		books = append(books, book)
	}
	c.IndentedJSON(http.StatusOK, books)
}

func HeadBooks(c *gin.Context) {
	FindBooks(c)
	c.Writer.WriteHeaderNow()
	c.Writer.Flush()
}

func Options(c *gin.Context) {
	c.Header("Allow", "GET, HEAD, POST, PATCH, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PATCH, DELETE, OPTIONS")
	c.Status(http.StatusOK)
}

func AddBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !models.ValidGenres[newBook.Genre] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre"})
		return
	}
	result, err := db.DB.Exec("INSERT INTO books (title, year, genre, status, link) VALUES (?, ?, ?, ?, ?)",
		newBook.Title, newBook.Year, newBook.Genre, newBook.Status, newBook.Link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new book"})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inserted ID"})
		return
	}
	newBook.ID = int(id)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func UpdateBook(c *gin.Context) {
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
	row := db.DB.QueryRow("SELECT title, year, genre, status, link FROM books WHERE id = ?", id)
	var currTitle string
	var currYear int
	var currGenre string
	var currStatus sql.NullBool
	var currLink sql.NullString
	if err := row.Scan(&currTitle, &currYear, &currGenre, &currStatus, &currLink); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	if update.Title != nil {
		currTitle = *update.Title
	}
	if update.Status != nil {
		currStatus = sql.NullBool{Bool: *update.Status, Valid: true}
	}
	if update.Year != nil {
		currYear = *update.Year
	}
	if update.Genre != nil {
		if !models.ValidGenres[*update.Genre] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre"})
			return
		}
		currGenre = *update.Genre
	}
	if update.Link != nil {
		currLink = sql.NullString{String: *update.Link, Valid: true}
	}
	_, err = db.DB.Exec(
		"UPDATE books SET title = ?, year = ?, genre = ?, status = ?, link = ? WHERE id = ?",
		currTitle, currYear, currGenre, currStatus.Bool, nullOrNil(currLink), id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book information"})
		return
	}
	resp := gin.H{
		"title":  currTitle,
		"year":   currYear,
		"genre":  currGenre,
		"status": currStatus.Bool,
	}
	if currLink.Valid {
		resp["link"] = currLink.String
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := db.DB.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check deletion result"})
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func nullOrNil(s sql.NullString) any {
	if s.Valid {
		return s.String
	}
	return nil
}
