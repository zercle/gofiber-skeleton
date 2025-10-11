package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/zercle/gofiber-skeleton/internal/config"
)

// Logger creates a logging middleware
func Logger(cfg *config.Config) fiber.Handler {
	logFormat := "[${time}] ${status} - ${method} ${path} - ${latency}\n"

	if cfg.Log.Format == "json" {
		logFormat = `{"time":"${time}","status":${status},"method":"${method}","path":"${path}","latency":"${latency}","ip":"${ip}","user_agent":"${ua}","request_id":"${locals:requestid}"}` + "\n"
	}

	return logger.New(logger.Config{
		Format:     logFormat,
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
		Output:     getLogOutput(cfg),
	})
}

// getLogOutput returns the log output based on configuration
func getLogOutput(cfg *config.Config) interface{} {
	// In production, you might want to use a file or external service
	// For now, we'll use stdout for all environments
	return nil // nil defaults to os.Stdout
}

// CustomLogger creates a custom structured logger
type CustomLogger struct {
	cfg *config.Config
}

// NewCustomLogger creates a new custom logger
func NewCustomLogger(cfg *config.Config) *CustomLogger {
	return &CustomLogger{cfg: cfg}
}

// Log logs a message with context
func (l *CustomLogger) Log(level, message string, fields map[string]interface{}) {
	if l.cfg.Log.Format == "json" {
		l.logJSON(level, message, fields)
	} else {
		l.logText(level, message, fields)
	}
}

// logJSON logs in JSON format
func (l *CustomLogger) logJSON(level, message string, fields map[string]interface{}) {
	fmt.Printf(`{"time":"%s","level":"%s","message":"%s"`, time.Now().Format(time.RFC3339), level, message)
	for key, value := range fields {
		fmt.Printf(`,"%s":"%v"`, key, value)
	}
	fmt.Println("}")
}

// logText logs in text format
func (l *CustomLogger) logText(level, message string, fields map[string]interface{}) {
	fmt.Printf("[%s] %s: %s", time.Now().Format(time.RFC3339), level, message)
	for key, value := range fields {
		fmt.Printf(" %s=%v", key, value)
	}
	fmt.Println()
}

// Info logs an info message
func (l *CustomLogger) Info(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log("INFO", message, f)
}

// Error logs an error message
func (l *CustomLogger) Error(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log("ERROR", message, f)
}

// Warn logs a warning message
func (l *CustomLogger) Warn(message string, fields ...map[string]interface{}) {
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log("WARN", message, f)
}

// Debug logs a debug message
func (l *CustomLogger) Debug(message string, fields ...map[string]interface{}) {
	if l.cfg.Log.Level != "debug" {
		return
	}
	f := make(map[string]interface{})
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log("DEBUG", message, f)
}
