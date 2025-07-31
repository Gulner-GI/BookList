package loggers

import (
	"log"
	"os"

	"github.com/Gulner-GI/BookList/config"
)

var LogFile *os.File

func InitLogger() {
	var err error
	LogFile, err = os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка при открытии лог-файла: %v", err)
	}
	log.SetOutput(LogFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Запуск сервера на %s", config.Port)
}
