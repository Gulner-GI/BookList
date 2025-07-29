package models

type Book struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Year       int     `json:"year"`
	Genre      string  `json:"genre"`
	Status     bool    `json:"-"`
	StatusText string  `json:"status"`
	Link       *string `json:"link,omitempty"`
}

type Update struct {
	Title  *string `json:"title,omitempty"`
	Year   *int    `json:"year,omitempty"`
	Genre  *string `json:"genre,omitempty"`
	Status *bool   `json:"status,omitempty"`
	Link   *string `json:"link,omitempty"`
}

var ValidGenres = map[string]bool{
	"Fantasy":     true,
	"Non-Fiction": true,
	"Sci-Fi":      true,
	"Science":     true,
}

type ErrorResponse400 struct {
	Message string `json:"message" example:"Некорректный ID"`
}
type ErrorResponse404 struct {
	Message string `json:"message" example:"Книга не найдена"`
}
type ErrorResponse500 struct {
	Message string `json:"message" example:"Внутренняя ошибка сервера"`
}
