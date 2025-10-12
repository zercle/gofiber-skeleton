package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	// Default logger instance
	defaultLogger *slog.Logger
)

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Level  string
	Format string
}

// NewLogger creates a new structured logger with the given configuration
func NewLogger(cfg LoggerConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	return logger
}

// Init initializes the default logger
func Init(cfg LoggerConfig) {
	defaultLogger = NewLogger(cfg)
}

// GetLogger returns the default logger instance
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		// Fallback to a basic logger if not initialized
		defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}
	return defaultLogger
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// WithContext returns a logger with context
func WithContext(ctx context.Context) *slog.Logger {
	return GetLogger()
}

// With returns a logger with additional attributes
func With(args ...any) *slog.Logger {
	return GetLogger().With(args...)
}