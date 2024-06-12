package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type CustomLogger struct {
	logFileName string
	stdLogger   *log.Logger
	doStdOut    bool
}

// Implement newCustomLogger
func NewCustomLogger(fn string) *CustomLogger {
	l := &CustomLogger{
		logFileName: fn,
		doStdOut:    true,
	}
	l.initStdLogger(fn)
	return l
}

func (l *CustomLogger) initStdLogger(fn string) {
	// Initializing the standard logger
	file := io.Writer(nil)
	err := error(nil)
	var stdOut io.Writer
	if l.doStdOut {
		stdOut = io.Writer(os.Stdout)
	} else {
		stdOut = file
	}
	if l.logFileName != "" {
		file, err = os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Failed to open log file: %v", err)
		}
	}
	l.stdLogger = log.New(io.MultiWriter(file, stdOut), "WHAT-TO: ", log.Ldate|log.Ltime)
}

func getLocationInfo() string {
	_, file, line, ok := runtime.Caller(2) // Caller(1) will return information about the caller of Info
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

// Add the Error method
func (l *CustomLogger) Error(str string, err error) {
	l.stdLogger.Printf("[ERROR] %s %s %v", getLocationInfo(), str, err)
}
