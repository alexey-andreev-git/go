package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type CustomLogger struct {
	logFileName string
	stdLogger   *log.Logger
}

// Implement newCustomLogger
func NewCustomLogger(fn string) *CustomLogger {
	l := &CustomLogger{
		logFileName: fn,
		stdLogger:   initStdLogger(fn),
	}
	return l
}

func initStdLogger(fn string) *log.Logger {
	// Initializing the standard logger
	file, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	return log.New(file, "WHAT-TO: ", log.Ldate|log.Ltime)
}

func getLocationInfo() string {
	_, file, line, ok := runtime.Caller(2) // Caller(1) вернет информацию о том, кто вызвал Info
	if !ok {
		file = "???"
		line = 0
	}
	logMessage := fmt.Sprintf("%s:%d:", filepath.Base(file), line)
	return logMessage
}

func (l *CustomLogger) Fatal(str string, err error) {
	l.stdLogger.Fatalf("[FATAL] %s %s %v", getLocationInfo(), str, err)
}

func (l *CustomLogger) Panic(str string, err error) {
	l.stdLogger.Panicf("[PANIC] %s %s %v", getLocationInfo(), str, err)
}

func (l *CustomLogger) Info(str string) {
	l.stdLogger.Printf("[INFO] %s %s", getLocationInfo(), str)
}

func (l *CustomLogger) Debug(str string) {
	l.stdLogger.Printf("[DEBUG] %s %s", getLocationInfo(), str)
}

func (l *CustomLogger) Warn(str string) {
	l.stdLogger.Printf("[WARN] %s %s", getLocationInfo(), str)
}
