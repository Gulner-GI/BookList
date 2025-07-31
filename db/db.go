package db

import (
	"database/sql"
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
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		year INTEGER NOT NULL,
		genre TEXT CHECK(genre IN ('Fantasy', 'Non-Fiction', 'Sci-Fi', 'Science')),
		status BOOLEAN NOT NULL CHECK (status IN (0,1)),
		link TEXT
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	log.Println("База данных подключена")
}
