package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Logger returns logging middleware
func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	})
}

// DetailedLogger returns detailed logging middleware for development
func DetailedLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path} | IP: ${ip} | UA: ${ua}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	})
}

// ProductionLogger returns structured JSON logging for production
func ProductionLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request
		duration := time.Since(start)

		// You can integrate with a structured logger here (e.g., zap, logrus)
		// For now, using simple format
		c.Context().Logger().Printf(
			`{"time":"%s","status":%d,"duration_ms":%d,"method":"%s","path":"%s","ip":"%s"}`,
			start.Format(time.RFC3339),
			c.Response().StatusCode(),
			duration.Milliseconds(),
			c.Method(),
			c.Path(),
			c.IP(),
		)

		return err
	}
}
