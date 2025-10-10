package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Config holds logger configuration
type Config struct {
	Environment string // production, development, or staging
	Level       string // debug, info, warn, error
	OutputPaths []string
}

// Init initializes the global logger
func Init(config Config) error {
	var zapConfig zap.Config

	// Set config based on environment
	switch config.Environment {
	case "production":
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	case "development":
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		zapConfig = zap.NewProductionConfig()
	}

	// Set log level
	if config.Level != "" {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(config.Level)); err == nil {
			zapConfig.Level = zap.NewAtomicLevelAt(level)
		}
	}

	// Set output paths
	if len(config.OutputPaths) > 0 {
		zapConfig.OutputPaths = config.OutputPaths
	} else {
		zapConfig.OutputPaths = []string{"stdout"}
	}

	var err error
	Log, err = zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	return nil
}

// InitDefault initializes logger with default configuration
func InitDefault(environment string) error {
	return Init(Config{
		Environment: environment,
		Level:       "info",
		OutputPaths: []string{"stdout"},
	})
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	}
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, err error, fields ...zap.Field) {
	if Log != nil {
		fields = append(fields, zap.Error(err))
		Log.Error(msg, fields...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	} else {
		// Fallback to stderr if logger not initialized
		os.Exit(1)
	}
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	if Log != nil {
		return Log.With(fields...)
	}
	return zap.NewNop()
}

// WithContext creates a logger with request context fields
func WithContext(requestID, userID, path string) *zap.Logger {
	fields := []zap.Field{
		zap.String("request_id", requestID),
		zap.String("path", path),
	}
	if userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	return With(fields...)
}

// LogHTTPRequest logs an HTTP request
func LogHTTPRequest(method, path, ip string, statusCode int, latency float64, requestID string) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.String("ip", ip),
		zap.Int("status", statusCode),
		zap.Float64("latency_ms", latency),
		zap.String("request_id", requestID),
	}

	if statusCode >= 500 {
		Error("HTTP request failed", nil, fields...)
	} else if statusCode >= 400 {
		Warn("HTTP request error", fields...)
	} else {
		Info("HTTP request", fields...)
	}
}

// LogDatabaseQuery logs a database query
func LogDatabaseQuery(query string, duration float64, err error) {
	fields := []zap.Field{
		zap.String("query", query),
		zap.Float64("duration_ms", duration),
	}

	if err != nil {
		Error("Database query failed", err, fields...)
	} else if duration > 1000 { // Slow query > 1s
		Warn("Slow database query", fields...)
	} else {
		Debug("Database query", fields...)
	}
}

// LogAuth logs authentication events
func LogAuth(action, userID, ip string, success bool, reason string) {
	fields := []zap.Field{
		zap.String("action", action),
		zap.String("user_id", userID),
		zap.String("ip", ip),
		zap.Bool("success", success),
	}

	if reason != "" {
		fields = append(fields, zap.String("reason", reason))
	}

	if success {
		Info("Authentication event", fields...)
	} else {
		Warn("Authentication failed", fields...)
	}
}

// LogBusinessEvent logs important business events
func LogBusinessEvent(event, userID string, metadata map[string]interface{}) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.String("user_id", userID),
		zap.Any("metadata", metadata),
	}
	Info("Business event", fields...)
}
