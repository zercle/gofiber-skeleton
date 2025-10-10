package response

import (
	"github.com/gofiber/fiber/v2"
)

// JSendStatus represents the status of a JSend response
type JSendStatus string

const (
	// StatusSuccess indicates successful request
	StatusSuccess JSendStatus = "success"
	// StatusFail indicates request failed due to client error
	StatusFail JSendStatus = "fail"
	// StatusError indicates request failed due to server error
	StatusError JSendStatus = "error"
)

// JSendResponse represents a JSend-compliant response
type JSendResponse struct {
	Status  JSendStatus `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// Success sends a success response with data
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(JSendResponse{
		Status: StatusSuccess,
		Data:   data,
	})
}

// Created sends a created response with data
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(JSendResponse{
		Status: StatusSuccess,
		Data:   data,
	})
}

// NoContent sends a no content response
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Fail sends a fail response (client error)
func Fail(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(JSendResponse{
		Status:  StatusFail,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response (server error)
func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(JSendResponse{
		Status:  StatusError,
		Message: message,
		Code:    status,
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *fiber.Ctx, message string, data interface{}) error {
	return Fail(c, fiber.StatusBadRequest, message, data)
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// Forbidden sends a 403 forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message)
}

// NotFound sends a 404 not found response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

// Conflict sends a 409 conflict response
func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, message)
}

// UnprocessableEntity sends a 422 unprocessable entity response
func UnprocessableEntity(c *fiber.Ctx, message string, data interface{}) error {
	return Fail(c, fiber.StatusUnprocessableEntity, message, data)
}

// InternalServerError sends a 500 internal server error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}

// ServiceUnavailable sends a 503 service unavailable response
func ServiceUnavailable(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusServiceUnavailable, message)
}
