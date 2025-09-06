package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omniti-labs/jsend"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
)

func SendSuccess(c *fiber.Ctx, data interface{}) error {
	return c.JSON(jsend.NewSuccess(data))
}

func SendCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(jsend.NewSuccess(data))
}

func SendError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(jsend.NewError(message))
}

func SendFail(c *fiber.Ctx, statusCode int, data map[string]interface{}) error {
	return c.Status(statusCode).JSON(jsend.NewFail(data))
}

func SendValidationError(c *fiber.Ctx, errors types.ValidationErrors) error {
	data := make(map[string]interface{})
	for _, err := range errors {
		data[err.Field] = err.Message
	}
	return SendFail(c, fiber.StatusBadRequest, data)
}

func SendNotFound(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return SendFail(c, fiber.StatusNotFound, map[string]interface{}{
		"message": message,
	})
}

func SendUnauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return SendFail(c, fiber.StatusUnauthorized, map[string]interface{}{
		"message": message,
	})
}

func SendForbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return SendFail(c, fiber.StatusForbidden, map[string]interface{}{
		"message": message,
	})
}

func SendConflict(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Conflict"
	}
	return SendFail(c, fiber.StatusConflict, map[string]interface{}{
		"message": message,
	})
}

func SendInternalError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return SendError(c, fiber.StatusInternalServerError, message)
}