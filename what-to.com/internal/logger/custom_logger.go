package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type customLogger struct {
	logFileName string
	stdLogger   *log.Logger
}

var (
	CustomLogger customLogger = newCustomLogger()
)

// Implement newCustomLogger
func newCustomLogger() customLogger {
	l := customLogger{}
	l.logFileName = "whattoapp.log"
	l.init()
	return l
}

func (l *customLogger) init() {
	// Initializing the standard logger
	file, err := os.OpenFile(l.logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	l.stdLogger = log.New(file, "WHAT-TO: ", log.Ldate|log.Ltime)
}

func (l *customLogger) Fatal(str string, err error) {
	l.stdLogger.Fatalf("[FATAL] %s %s %v", getLocationInfo(), str, err)
}

func (l *customLogger) Panic(str string, err error) {
	l.stdLogger.Panicf("[PANIC] %s %s %v", getLocationInfo(), str, err)
}

func (l *customLogger) Info(str string) {
	l.stdLogger.Printf("[INFO] %s %s", getLocationInfo(), str)
}

func (l *customLogger) Debug(str string) {
	l.stdLogger.Printf("[DEBUG] %s %s", getLocationInfo(), str)
}

func (l *customLogger) Warn(str string) {
	l.stdLogger.Printf("[WARN] %s %s", getLocationInfo(), str)
}

func getLocationInfo() string {
	_, file, line, ok := runtime.Caller(2) // Caller(1) вернет информацию о том, кто вызвал Info
	if !ok {
		file = "???"
		line = 0
	}
	logMessage := fmt.Sprintf("%s:%d:", filepath.Base(file), line)
	// l.stdLogger.Printf(logMessage)
	return logMessage
}
