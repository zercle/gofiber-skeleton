package middleware

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/zercle/gofiber-skeleton/internal/errors"
	"github.com/zercle/gofiber-skeleton/internal/response"
)

// ErrorHandler returns a centralized error handling middleware
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process the request
		err := c.Next()

		// If no error, continue
		if err == nil {
			return nil
		}

		// Log the error (in production, use structured logging)
		requestID := c.Locals("requestid")
		log.Printf("[ERROR] RequestID: %v, Path: %s, Error: %v", requestID, c.Path(), err)

		// Handle domain errors
		if domainErr, ok := err.(*errors.DomainError); ok {
			return handleDomainError(c, domainErr)
		}

		// Handle Fiber errors
		if fiberErr, ok := err.(*fiber.Error); ok {
			return handleFiberError(c, fiberErr)
		}

		// Handle PostgreSQL errors
		if pqErr, ok := err.(*pq.Error); ok {
			return handlePostgresError(c, pqErr)
		}

		// Handle SQL errors
		if err == sql.ErrNoRows {
			return response.NotFound(c, "Record not found")
		}

		// Default to internal server error
		return response.InternalServerError(c, "An internal error occurred")
	}
}

// handleDomainError handles custom domain errors
func handleDomainError(c *fiber.Ctx, err *errors.DomainError) error {
	statusCode := err.HTTPStatus
	if statusCode == 0 {
		statusCode = fiber.StatusInternalServerError
	}

	// Build response
	resp := map[string]interface{}{
		"status":  "error",
		"message": err.Message,
		"code":    err.Code,
	}

	// Add context in development mode (don't expose in production)
	if c.Locals("env") == "development" && len(err.Context) > 0 {
		resp["context"] = err.Context
	}

	return c.Status(statusCode).JSON(resp)
}

// handleFiberError handles Fiber framework errors
func handleFiberError(c *fiber.Ctx, err *fiber.Error) error {
	code := err.Code
	if code == 0 {
		code = fiber.StatusInternalServerError
	}

	return response.Error(c, code, err.Message)
}

// handlePostgresError handles PostgreSQL-specific errors
func handlePostgresError(c *fiber.Ctx, err *pq.Error) error {
	switch err.Code {
	case "23505": // unique_violation
		return response.Conflict(c, "Duplicate record: "+err.Detail)

	case "23503": // foreign_key_violation
		return response.BadRequest(c, "Foreign key constraint violation", nil)

	case "23502": // not_null_violation
		return response.BadRequest(c, fmt.Sprintf("Required field missing: %s", err.Column), nil)

	case "22P02": // invalid_text_representation
		return response.BadRequest(c, "Invalid data format", nil)

	case "42P01": // undefined_table
		return response.InternalServerError(c, "Database configuration error")

	default:
		// Log detailed error for debugging
		log.Printf("[PG_ERROR] Code: %s, Message: %s, Detail: %s", err.Code, err.Message, err.Detail)

		return response.InternalServerError(c, "Database error occurred")
	}
}

// NotFoundHandler returns a custom 404 handler
func NotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.NotFound(c, fmt.Sprintf("Route '%s' not found", c.Path()))
	}
}

// MethodNotAllowedHandler returns a custom 405 handler
func MethodNotAllowedHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusMethodNotAllowed,
			fmt.Sprintf("Method '%s' not allowed for '%s'", c.Method(), c.Path()))
	}
}
