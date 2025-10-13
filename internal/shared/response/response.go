package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Data      any       `json:"data,omitempty"`
	Error     any       `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func OK(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func NoContent(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNoContent).JSON(Response{
		Success:   true,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func BadRequest(c *fiber.Ctx, message string, details any) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success:   false,
		Message:   message,
		Error:     details,
		Timestamp: time.Now(),
	})
}

func ValidationError(c *fiber.Ctx, details any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Success:   false,
		Message:   "Validation error",
		Error:     details,
		Timestamp: time.Now(),
	})
}

func Unauthorized(c *fiber.Ctx, message string, details any) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success:   false,
		Message:   message,
		Error:     details,
		Timestamp: time.Now(),
	})
}

func Forbidden(c *fiber.Ctx, message string, details any) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success:   false,
		Message:   message,
		Error:     details,
		Timestamp: time.Now(),
	})
}

func NotFound(c *fiber.Ctx, message string, details any) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success:   false,
		Message:   message,
		Error:     details,
		Timestamp: time.Now(),
	})
}

func Conflict(c *fiber.Ctx, message string, details any) error {
	return c.Status(fiber.StatusConflict).JSON(Response{
		Success:   false,
		Message:   message,
		Error:     details,
		Timestamp: time.Now(),
	})
}

func InternalServerError(c *fiber.Ctx, message string, details any) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success:   false,
		Message:   message,
		Error:     details,
		Timestamp: time.Now(),
	})
}
