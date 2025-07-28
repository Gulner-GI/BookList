package main

import (
	"github.com/Gulner-GI/ToDoApp/db"
	"github.com/Gulner-GI/ToDoApp/routes"
)

func main() {
	db.InitDB("todo.db")
	r := routes.SetupRouter()
	r.Run(":8080")
}
