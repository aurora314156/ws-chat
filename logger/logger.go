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
func Info(v ...any) {
	log.Println("[ℹ️ INFO]", v)
}

// Error log
func Error(v ...any) {
	log.Println("[❌ ERROR]", v)
}

// Debug log
func Debug(v ...any) {
	log.Println("[🐛 DEBUG]", v)
}
