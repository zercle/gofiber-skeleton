package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// Recovery creates a recovery middleware
func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				requestID := GetRequestID(c)
				stackTrace := debug.Stack()

				errorMessage := fmt.Sprintf("[PANIC] RequestID: %s | Path: %s | Error: %v\nStack Trace:\n%s",
					requestID,
					c.Path(),
					r,
					string(stackTrace),
				)
				fmt.Println(errorMessage)

				_ = c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
					Error:   "Internal Server Error",
					Message: "An unexpected error occurred",
					Code:    "INTERNAL_ERROR",
				})
			}
		}()

		return c.Next()
	}
}

// ErrorHandler creates a centralized error handler
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process request
		err := c.Next()

		// If no error, continue
		if err == nil {
			return nil
		}

		// Check if it's a fiber.Error
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(ErrorResponse{
				Error: fiberErr.Message,
			})
		}

		// Handle other errors
		requestID := GetRequestID(c)
		errorMessage := fmt.Sprintf("[%s] ERROR | RequestID: %s | Path: %s | Error: %v",
			"ERROR",
			requestID,
			c.Path(),
			err,
		)
		fmt.Println(errorMessage)

		// Return generic error to client
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: "An unexpected error occurred",
			Code:    "INTERNAL_ERROR",
		})
	}
}

// NotFoundHandler creates a custom 404 handler
func NotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Not Found",
			Message: fmt.Sprintf("Route %s %s not found", c.Method(), c.Path()),
			Code:    "NOT_FOUND",
		})
	}
}

// MethodNotAllowedHandler creates a custom 405 handler
func MethodNotAllowedHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(ErrorResponse{
			Error:   "Method Not Allowed",
			Message: fmt.Sprintf("Method %s not allowed for route %s", c.Method(), c.Path()),
			Code:    "METHOD_NOT_ALLOWED",
		})
	}
}

// RateLimitHandler creates a rate limit response handler
func RateLimitHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(ErrorResponse{
			Error:   "Rate Limit Exceeded",
			Message: "Too many requests, please try again later",
			Code:    "RATE_LIMIT_EXCEEDED",
		})
	}
}
