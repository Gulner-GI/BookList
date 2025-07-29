// @title BookList API
// @version 0.3
// @description RESTful API on gin and SQLite3 for listing books
// @termsOfService http://swagger.io/terms/

// @contact.name Gulner
// @contact.email yura.kuzin1990@gmail.com

// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"github.com/Gulner-GI/BookList/config"
	"github.com/Gulner-GI/BookList/db"
	"github.com/Gulner-GI/BookList/loggers"
	"github.com/Gulner-GI/BookList/routes"
)

func main() {
	loggers.InitLogger()
	db.InitDB(config.DBPath)
	r := routes.SetupRouter()
	r.Run(config.Port)
}
