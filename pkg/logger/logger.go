package logger

import (
	"fmt"
	"log"
	"whatsapp-api-go/internal/domain/ports"
)

// SimpleLogger implementa ports.Logger
type SimpleLogger struct{}

// NewSimpleLogger crea un nuevo logger simple
func NewSimpleLogger() ports.Logger {
	return &SimpleLogger{}
}

// Info registra un mensaje informativo
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		log.Printf("[INFO] %s %v\n", msg, args)
	} else {
		log.Printf("[INFO] %s\n", msg)
	}
}

// Error registra un mensaje de error
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		log.Printf("[ERROR] %s %v\n", msg, args)
	} else {
		log.Printf("[ERROR] %s\n", msg)
	}
}

// Debug registra un mensaje de depuración
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		log.Printf("[DEBUG] %s %v\n", msg, args)
	} else {
		log.Printf("[DEBUG] %s\n", msg)
	}
}

// Warn registra un mensaje de advertencia
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		log.Printf("[WARN] %s %v\n", msg, args)
	} else {
		log.Printf("[WARN] %s\n", msg)
	}
}

// ColorLogger implementa ports.Logger con colores
type ColorLogger struct{}

// NewColorLogger crea un nuevo logger con colores
func NewColorLogger() ports.Logger {
	return &ColorLogger{}
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

func (l *ColorLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf("%s[INFO]%s %s %v\n", colorGreen, colorReset, msg, args)
	} else {
		fmt.Printf("%s[INFO]%s %s\n", colorGreen, colorReset, msg)
	}
}

func (l *ColorLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf("%s[ERROR]%s %s %v\n", colorRed, colorReset, msg, args)
	} else {
		fmt.Printf("%s[ERROR]%s %s\n", colorRed, colorReset, msg)
	}
}

func (l *ColorLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf("%s[DEBUG]%s %s %v\n", colorGray, colorReset, msg, args)
	} else {
		fmt.Printf("%s[DEBUG]%s %s\n", colorGray, colorReset, msg)
	}
}

func (l *ColorLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf("%s[WARN]%s %s %v\n", colorYellow, colorReset, msg, args)
	} else {
		fmt.Printf("%s[WARN]%s %s\n", colorYellow, colorReset, msg)
	}
}

