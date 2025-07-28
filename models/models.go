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
	Title  *string `json:"task"`
	Year   *int    `json:"year"`
	Genre  *string `json:"genre"`
	Status *bool   `json:"status"`
	Link   *string `json:"link,omitempty"`
}

var ValidGenres = map[string]bool{
	"Fantasy":     true,
	"Non-Fiction": true,
	"Sci-Fi":      true,
	"Science":     true,
}
