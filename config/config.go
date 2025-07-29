package config

import "os"

func mustGet(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

var (
	Port   = ":" + mustGet("PORT", "8080")
	DBPath = mustGet("DB_PATH", "books.db")
)
