package utils

import (
	"log"
	"os"
	"runtime"
	"time"

	"health-store/config"
)

// LogLevel represents the level of logging
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger wraps the standard log package with additional functionality
type Logger struct {
	*log.Logger
	level LogLevel
	env   string
}

// Global logger instance
var AppLogger *Logger

// InitLogger initializes the global logger
func InitLogger(cfg *config.Config) {
	AppLogger = &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		level:  getLogLevel(cfg.Server.Env),
		env:    cfg.Server.Env,
	}
}

// getLogLevel returns the appropriate log level based on environment
func getLogLevel(env string) LogLevel {
	switch env {
	case "production":
		return WARN
	case "development":
		return DEBUG
	case "test":
		return ERROR
	default:
		return INFO
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	if l.level <= DEBUG {
		l.Printf("[DEBUG] %s", message)
	}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	if l.level <= INFO {
		l.Printf("[INFO] %s", message)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string) {
	if l.level <= WARN {
		l.Printf("[WARN] %s", message)
	}
}

// Error logs an error message
func (l *Logger) Error(message string) {
	if l.level <= ERROR {
		l.Printf("[ERROR] %s", message)
	}
}

// Debugf logs a debug message with formatting
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.Printf("[DEBUG] "+format, v...)
	}
}

// Infof logs an info message with formatting
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= INFO {
		l.Printf("[INFO] "+format, v...)
	}
}

// Warnf logs a warning message with formatting
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= WARN {
		l.Printf("[WARN] "+format, v...)
	}
}

// Errorf logs an error message with formatting
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.Printf("[ERROR] "+format, v...)
	}
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(method, path, userAgent, ip string, statusCode int, duration time.Duration) {
	if l.level <= INFO {
		l.Infof("HTTP %s %s | %d | %v | %s | %s", method, path, statusCode, duration, ip, userAgent)
	}
}

// LogError logs errors with context
func (l *Logger) LogError(err error, message string) {
	if l.level <= ERROR {
		_, file, line, _ := runtime.Caller(1)
		l.Errorf("Error in %s:%d - %s: %v", file, line, message, err)
	}
}

// Convenience functions for global logger
func Debug(message string) {
	if AppLogger != nil {
		AppLogger.Debug(message)
	}
}

func Info(message string) {
	if AppLogger != nil {
		AppLogger.Info(message)
	}
}

func Warn(message string) {
	if AppLogger != nil {
		AppLogger.Warn(message)
	}
}

func Error(message string) {
	if AppLogger != nil {
		AppLogger.Error(message)
	}
}

func Debugf(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debugf(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Infof(format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Warnf(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Errorf(format, v...)
	}
}

func LogError(err error, message string) {
	if AppLogger != nil {
		AppLogger.LogError(err, message)
	}
}
