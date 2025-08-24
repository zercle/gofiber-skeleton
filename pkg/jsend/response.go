package jsend

import (
	"github.com/gofiber/fiber/v2"
)

var Empty struct{}

// JSendResponse represents the base JSend response structure
// T is the type of the Data field.
type JSendResponse[T any] struct {
	Status  string `json:"status"`                  // "success", "fail", or "error"
	Data    T      `json:"data,omitempty,omitzero"` // Data for success/fail, omitted for error (if zero value is omittable)
	Message string `json:"message,omitempty"`       // Message for fail/error, omitted for success
	Code    int    `json:"code,omitempty"`          // Optional error code for error status
}

// ErrorResponse represents a JSend error response with no data field
type ErrorResponse struct {
	Status  string `json:"status"`  // "error"
	Message string `json:"message"` // Message for error
	Code    int    `json:"code"`    // Error code
}

// Success sends a JSend success response
func Success[T any](c *fiber.Ctx, data T) error {
	return c.Status(fiber.StatusOK).JSON(JSendResponse[T]{
		Status: "success",
		Data:   data,
	})
}

// SuccessWithStatus sends a JSend success response with a custom HTTP status code
func SuccessWithStatus[T any](c *fiber.Ctx, data T, status int) error {
	return c.Status(status).JSON(JSendResponse[T]{
		Status: "success",
		Data:   data,
	})
}

// Fail sends a JSend fail response
func Fail[T any](c *fiber.Ctx, data T, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(JSendResponse[T]{
		Status:  "fail",
		Data:    data,
		Message: message,
	})
}

// Error sends a JSend error response
func Error(c *fiber.Ctx, message string, code int, status int) error {
	return c.Status(status).JSON(JSendResponse[any]{
		Status:  "error",
		Message: message,
		Code:    code,
		Data:    nil, // Data is typically null or omitted for error responses
	})
}
