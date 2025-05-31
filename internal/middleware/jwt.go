package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gofiber-skeleton/pkg/jwtutil"
)

func JWTMiddleware(jwt *jwtutil.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or missing authorization header"})
		}
		token := parts[1]
		if !jwt.ValidateToken(token) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	}
}