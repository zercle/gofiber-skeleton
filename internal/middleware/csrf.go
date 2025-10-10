package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

// CSRF returns CSRF protection middleware with default configuration
func CSRF(secret string) fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		Expiration:     1 * time.Hour,
		KeyGenerator:   func() string { return secret },
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "CSRF token validation failed",
			})
		},
	})
}

// CSRFWithConfig returns CSRF protection middleware with custom configuration
func CSRFWithConfig(config csrf.Config) fiber.Handler {
	// Set sensible defaults if not provided
	if config.KeyLookup == "" {
		config.KeyLookup = "header:X-CSRF-Token"
	}
	if config.CookieName == "" {
		config.CookieName = "csrf_"
	}
	if config.Expiration == 0 {
		config.Expiration = 1 * time.Hour
	}
	if config.CookieSameSite == "" {
		config.CookieSameSite = "Strict"
	}

	return csrf.New(config)
}

// GetCSRFToken returns a handler that provides the CSRF token to clients
func GetCSRFToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("csrf").(string)
		return c.JSON(fiber.Map{
			"csrf_token": token,
		})
	}
}
