package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/pkg/response"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		log.Printf("Error: %v", err)

		var resp *response.JSendResponse
		
		switch code {
		case fiber.StatusBadRequest:
			resp = response.Fail(map[string]interface{}{
				"error": "Bad request",
			})
		case fiber.StatusUnauthorized:
			resp = response.Fail(map[string]interface{}{
				"error": "Unauthorized",
			})
		case fiber.StatusForbidden:
			resp = response.Fail(map[string]interface{}{
				"error": "Forbidden",
			})
		case fiber.StatusNotFound:
			resp = response.Fail(map[string]interface{}{
				"error": "Resource not found",
			})
		case fiber.StatusUnprocessableEntity:
			resp = response.Fail(map[string]interface{}{
				"error": "Validation failed",
			})
		default:
			resp = response.Error("Internal server error")
		}

		return c.Status(code).JSON(resp)
	}
}