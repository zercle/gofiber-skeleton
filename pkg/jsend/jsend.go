package jsend

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type JSendResponse struct {
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func Success(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusOK).JSON(JSendResponse{
		Status: "success",
		Data:   data,
	})
}

func Fail(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusBadRequest).JSON(JSendResponse{
		Status: "fail",
		Data:   data,
	})
}

func Error(c *fiber.Ctx, message string, code int, statusCode ...int) error {
	status := http.StatusInternalServerError
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	return c.Status(status).JSON(JSendResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	})
}
