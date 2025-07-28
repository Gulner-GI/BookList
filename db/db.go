package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		done BOOLEAN NOT NULL CHECK (done IN (0,1))
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	fmt.Println("База данных подключена и таблица создана (если её не было)")
}
