package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// RequestID creates a request ID middleware
func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header:     "X-Request-ID",
		ContextKey: "requestid",
	})
}
