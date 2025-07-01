package jsend

import (
	"github.com/gofiber/fiber/v2"
)

type JSendResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(JSendResponse{
		Status: "success",
		Data:    data,
	})
}

func Fail(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(JSendResponse{
		Status: "fail",
		Data:    data,
	})
}

func Error(c *fiber.Ctx, message string, code int) error {
	return c.Status(fiber.StatusInternalServerError).JSON(JSendResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	})
}