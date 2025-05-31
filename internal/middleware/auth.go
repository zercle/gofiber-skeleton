package middleware

import (
	"gofiber-skeleton/pkg/jwtutil"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const UserIDContextKey = "user_id"

func AuthMiddleware(jwt *jwtutil.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header"})
		}
		token := parts[1]
		if !jwt.ValidateToken(token) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		userID, err := jwt.ExtractUserID(token)
		if err != nil || userID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user in token"})
		}
		// Set user ID in context for downstream handlers
		c.Locals(UserIDContextKey, strconv.FormatUint(uint64(userID), 10))
		return c.Next()
	}
}
