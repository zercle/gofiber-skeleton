package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/template-go-fiber/internal/errors"
)

// Response is a generic API response wrapper
type Response[T any] struct {
	Success   bool   `json:"success"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Data      T      `json:"data,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// ErrorResponse is an error API response wrapper
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// SendSuccess sends a successful response
func SendSuccess[T any](c *fiber.Ctx, statusCode int, data T) error {
	return c.Status(statusCode).JSON(Response[T]{
		Success:   true,
		Data:      data,
		Timestamp: getCurrentUnix(),
	})
}

// SendCreated sends a 201 Created response
func SendCreated[T any](c *fiber.Ctx, data T) error {
	return SendSuccess(c, fiber.StatusCreated, data)
}

// SendOK sends a 200 OK response
func SendOK[T any](c *fiber.Ctx, data T) error {
	return SendSuccess(c, fiber.StatusOK, data)
}

// SendNoContent sends a 204 No Content response
func SendNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// SendError sends an error response
func SendError(c *fiber.Ctx, err *errors.APIError) error {
	return c.Status(err.StatusCode).JSON(ErrorResponse{
		Success:   false,
		Code:      string(err.Code),
		Message:   err.Message,
		Timestamp: getCurrentUnix(),
	})
}

// SendUnknownError sends a generic internal server error response
func SendUnknownError(c *fiber.Ctx, err error) error {
	apiErr := errors.AsAPIError(err)
	return SendError(c, apiErr)
}

// Helper function to get current Unix timestamp
func getCurrentUnix() int64 {
	return time.Now().Unix()
}
