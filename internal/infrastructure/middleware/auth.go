package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/omniti-labs/jsend"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

type AuthClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(jsend.NewFail(map[string]interface{}{
				"message": "Authorization header is required",
			}))
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(jsend.NewFail(map[string]interface{}{
				"message": "Invalid authorization header format",
			}))
		}

		tokenString := parts[1]
		claims := &AuthClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(jsend.NewFail(map[string]interface{}{
				"message": "Invalid or expired token",
			}))
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("claims", claims)

		return c.Next()
	}
}

func OptionalAuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]
		claims := &AuthClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err == nil && token.Valid {
			c.Locals("userID", claims.UserID)
			c.Locals("userEmail", claims.Email)
			c.Locals("claims", claims)
		}

		return c.Next()
	}
}