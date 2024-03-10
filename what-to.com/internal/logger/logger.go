package logger

import (
	"log"
	"os"
)

var (
	// Log - экспортированный логгер для использования в других пакетах
	Log *log.Logger
)

func init() {
	// Инициализация стандартного логгера
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	Log = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
