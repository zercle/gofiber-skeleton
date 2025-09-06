package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/omniti-labs/jsend"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

func NewRecover(cfg *config.Config) fiber.Handler {
	return recover.New(recover.Config{
		Next:              nil,
		EnableStackTrace:  cfg.IsDevelopment(),
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("Panic recovered: %v", e)
		},
	})
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("Error: %v", err)

	if code >= 500 {
		return c.Status(code).JSON(jsend.NewError(message))
	}

	return c.Status(code).JSON(jsend.NewFail(map[string]interface{}{
		"message": message,
	}))
}