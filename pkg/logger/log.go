package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLogger    = log.New(os.Stdout, "\033[36mℹ️ INFO:\033[0m ", 0)
	warnLogger    = log.New(os.Stdout, "\033[33m⚠️ WARN:\033[0m ", 0)
	errorLogger   = log.New(os.Stderr, "\033[31m❌ ERROR:\033[0m ", 0)
	successLogger = log.New(os.Stdout, "\033[32m✅ SUCCESS:\033[0m ", 0)
)

// Info logs an informational message
func Info(msg string) {
	coloredMsg := fmt.Sprintf("\033[36m%s\033[0m", msg) // Cyan
	infoLogger.Println(coloredMsg)
}

// Warn logs a warning message
func Warn(msg string) {
	coloredMsg := fmt.Sprintf("\033[33m%s\033[0m", msg) // Yellow
	warnLogger.Println(coloredMsg)
}

// Error logs an error
func Error(err error) {
	coloredMsg := fmt.Sprintf("\033[31m%s\033[0m", err.Error()) // Red
	errorLogger.Println(coloredMsg)
}

// Success logs a success message
func Success(msg string) {
	coloredMsg := fmt.Sprintf("\033[32m%s\033[0m", msg) // Green
	successLogger.Println(coloredMsg)
}
