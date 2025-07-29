package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Gulner-GI/BookList/db"
	"github.com/Gulner-GI/BookList/models"

	"github.com/gin-gonic/gin"
)

// FindBooks godoc
// @Summary Получение книг
// @Description Получает список всех книг, либо конкретную книгу по ID
// @Tags books
// @Accept json
// @Produce json
// @Param id query int false "ID книги для поиска"
// @Success 200 {object} models.Book     "Если указан ID — возвращается одна книга"
// @Success 200 {array}  models.Book     "Если ID не указан — возвращается список всех книг"
// @Failure 400 {object} models.ErrorResponse400 "Некорректный ID"
// @Failure 404 {object} models.ErrorResponse404 "Книга не найдена"
// @Failure 500 {object} models.ErrorResponse500 "Внутренняя ошибка сервера"
// @Router /books [get]
// @Router /books [head]
func FindBooks(c *gin.Context) {
	idParam := c.Query("id")
	if c.Request.Method == http.MethodHead {
		log.Println("HEAD-запрос к /books")
		c.Status(http.StatusOK)
		return
	}
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Printf("Некорректный ID в запросе: %v", idParam)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		log.Printf("Поиск книги с ID: %d", id)
		row := db.DB.QueryRow("SELECT id, title, year, genre, status, link FROM books WHERE id = ?", id)
		var book models.Book
		err = row.Scan(&book.ID, &book.Title, &book.Year, &book.Genre, &book.Status, &book.Link)
		if err == sql.ErrNoRows {
			log.Printf("Книга с ID %d не найдена", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		} else if err != nil {
			log.Printf("Ошибка при запросе книги с ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if book.Status {
			book.StatusText = "completed"
		} else {
			book.StatusText = "in process"
		}
		log.Printf("Книга найдена: %+v", book)
		c.IndentedJSON(http.StatusOK, book)
		return
	}

	log.Println("Запрос на получение всех книг")
	rows, err := db.DB.Query("SELECT id, title, year, genre, status, link FROM books")
	if err != nil {
		log.Printf("Ошибка при получении всех книг: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Year, &book.Genre, &book.Status, &book.Link)
		if err != nil {
			log.Printf("Ошибка при сканировании книги: %v", err)
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
	log.Printf("Получено книг: %d", len(books))
	c.IndentedJSON(http.StatusOK, books)
}

// HeadBooks godoc
// @Summary      HEAD-запрос для проверки доступности ресурса
// @Description  Возвращает только заголовки, без тела. Используется для проверки доступности /books
// @Tags         books
// @Produce      json
// @Success      200 {string} string "OK"
// @Router       /books [head]
func HeadBooks(c *gin.Context) {
	log.Println("HEAD-запрос к /books")
	FindBooks(c)
	c.Writer.WriteHeaderNow()
	c.Writer.Flush()
}

// Options godoc
// @Summary      OPTIONS-запрос к /books
// @Description  Возвращает список доступных методов
// @Tags         books
// @Produce      json
// @Success      200 {string} string "OK"
// @Router       /books [options]
func Options(c *gin.Context) {
	log.Println("OPTIONS-запрос к /books")
	c.Header("Allow", "GET, HEAD, POST, PATCH, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PATCH, DELETE, OPTIONS")
	c.Status(http.StatusOK)
}

// AddBook godoc
// @Summary      Добавить новую книгу
// @Description  Добавляет книгу в базу данных. Поле link опционально, можно удалить
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book body models.Book true "Данные новой книги"
// @Success      201 {object} models.Book
// @Failure 400 {object} models.ErrorResponse400 "Некорректный ID"
// @Failure 500 {object} models.ErrorResponse500 "Внутренняя ошибка сервера"
// @Router       /books [post]
func AddBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		log.Printf("Ошибка парсинга JSON для новой книги: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !models.ValidGenres[newBook.Genre] {
		log.Printf("Попытка добавить книгу с некорректным жанром: %s", newBook.Genre)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre"})
		return
	}
	result, err := db.DB.Exec("INSERT INTO books (title, year, genre, status, link) VALUES (?, ?, ?, ?, ?)",
		newBook.Title, newBook.Year, newBook.Genre, newBook.Status, newBook.Link)
	if err != nil {
		log.Printf("Ошибка вставки книги в базу: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new book"})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Ошибка получения ID вставленной книги: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inserted ID"})
		return
	}
	newBook.ID = int(id)
	log.Printf("Книга успешно добавлена с ID %d", newBook.ID)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// UpdateBook обновляет данные о книге по ID.
// @Summary      Обновить книгу
// @Description  Обновляет информацию о книге. Все поля опциональны (удалите нежелательные). Жанр должен быть из допустимого списка: Fantasy, Sci-FI, Science, Non-Fiction
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int                  true  "ID книги"
// @Param        data  body      models.Update        true  "Данные для обновления книги"
// @Success      200   {object}  models.Book          "Обновлённая книга"
// @Failure		 400 {object} models.ErrorResponse400 "Некорректный ID"
// @Failure		 404 {object} models.ErrorResponse404 "Книга не найдена"
// @Failure 	 500 {object} models.ErrorResponse500 "Внутренняя ошибка сервера"
// @Router       /books/{id} [patch]
func UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Некорректный ID для обновления: %v", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var update models.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		log.Printf("Ошибка парсинга JSON для обновления книги ID %d: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Обновление книги с ID %d, входящие данные: %+v", id, update)
	row := db.DB.QueryRow("SELECT title, year, genre, status, link FROM books WHERE id = ?", id)
	var currTitle string
	var currYear int
	var currGenre string
	var currStatus sql.NullBool
	var currLink sql.NullString
	if err := row.Scan(&currTitle, &currYear, &currGenre, &currStatus, &currLink); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Книга с ID %d не найдена для обновления", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			log.Printf("Ошибка при запросе книги с ID %d: %v", id, err)
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
			log.Printf("Попытка обновить книгу ID %d с некорректным жанром: %s", id, *update.Genre)
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
		log.Printf("Ошибка при обновлении книги ID %d: %v", id, err)
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
	log.Printf("Книга ID %d успешно обновлена: %+v", id, resp)
	c.IndentedJSON(http.StatusOK, resp)
}

// DeleteBook удаляет книгу по ID.
// @Summary      Удалить книгу
// @Description  Удаляет книгу из базы данных по её ID
// @Tags         books
// @Param        id  path  int  true  "ID книги"
// @Success      204  "Книга удалена"
// @Failure 400 {object} models.ErrorResponse400 "Некорректный ID"
// @Failure 404 {object} models.ErrorResponse404 "Книга не найдена"
// @Failure 500 {object} models.ErrorResponse500 "Внутренняя ошибка сервера"
// @Router       /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Некорректный ID для удаления: %v", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	log.Printf("Запрос на удаление книги с ID: %d", id)
	result, err := db.DB.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		log.Printf("Ошибка при удалении книги с ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Ошибка при проверке результата удаления книги с ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check deletion result"})
		return
	}
	if rowsAffected == 0 {
		log.Printf("Книга с ID %d не найдена для удаления", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	log.Printf("Книга с ID %d успешно удалена", id)
	c.Status(http.StatusNoContent)
}

func nullOrNil(s sql.NullString) any {
	if s.Valid {
		return s.String
	}
	return nil
}
