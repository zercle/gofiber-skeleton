package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Logger creates a simple logger middleware
func Logger() fiber.Handler {
	return CustomLogger()
}

// RequestID adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate unique request ID
		requestID := uuid.New().String()

		// Store in context
		c.Locals("request_id", requestID)

		// Add to response headers
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(c *fiber.Ctx) string {
	if requestID, ok := c.Locals("request_id").(string); ok {
		return requestID
	}
	return ""
}

// CustomLogger creates a custom logger that includes request ID
func CustomLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID
		requestID := GetRequestID(c)

		// Log the request
		logMessage := fmt.Sprintf("%s | %d | %v | %s | %s | %s | %s",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Response().StatusCode(),
			latency,
			c.Method(),
			c.Path(),
			c.IP(),
			requestID,
		)

		if err != nil {
			logMessage += fmt.Sprintf(" | %s", err.Error())
		}

		fmt.Println(logMessage)

		return err
	}
}

// ErrorLogger creates an error-specific logger
func ErrorLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			requestID := GetRequestID(c)
			errorMessage := fmt.Sprintf("[%s] ERROR | RequestID: %s | Method: %s | Path: %s | Error: %v",
				time.Now().Format("2006-01-02 15:04:05"),
				requestID,
				c.Method(),
				c.Path(),
				err,
			)
			fmt.Println(errorMessage)
		}

		return err
	}
}