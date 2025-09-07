package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/shared/jsend"
)

func NewRecover(cfg *config.Config) RecoverMiddleware {
	return RecoverMiddleware(recover.New(recover.Config{
		EnableStackTrace: cfg.IsDevelopment(),
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("Panic occurred: %v", e)
			if cfg.IsDevelopment() {
				jsend.SendError(c, fiber.StatusInternalServerError, "Internal Server Error: Panic Recovered", 0)
			} else {
				jsend.SendError(c, fiber.StatusInternalServerError, "Internal Server Error", 0)
			}
		},
		Next: nil, // Let the next middleware or route handle the request, or the global error handler
	}))
}
