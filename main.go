package main

import (
	"github.com/Gulner-GI/ToDoApp/routes"
)

func main() {
	r := routes.SetupRouter()
	r.Run(":8080")
}
