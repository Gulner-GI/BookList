package main

import (
	"github.com/Gulner-GI/BookList/db"
	"github.com/Gulner-GI/BookList/routes"
)

func main() {
	db.InitDB("books.db")
	r := routes.SetupRouter()
	r.Run(":8080")
}
