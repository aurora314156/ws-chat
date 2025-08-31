package logger

import (
	"log"
	"os"
)

// Log levels
type Level int

const (
	INFO Level = iota
	ERROR
	DEBUG
)

func init() {
	// Initialize the logger to write to standard output with date, time, and file info
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Info log
func Info(format string, v ...any) {
	log.Printf("[ℹ️ INFO] "+format, v...)
}

// Error log
func Error(format string, v ...any) {
	log.Printf("[❌ ERROR]"+format, v)
}

// Debug log
func Debug(format string, v ...any) {
	log.Printf("[🐛 DEBUG]"+format, v)
}
