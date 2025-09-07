package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/shared/jsend"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuth(cfg *config.Config) AuthMiddleware {
	return AuthMiddleware(func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return jsend.SendUnauthorized(c, map[string]string{
				"message": "Missing authorization header",
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return jsend.SendUnauthorized(c, map[string]string{
				"message": "Invalid authorization header format",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return jsend.SendUnauthorized(c, map[string]string{
				"message": "Missing token",
			})
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			return jsend.SendUnauthorized(c, map[string]string{
				"message": "Invalid token",
				"error":   err.Error(),
			})
		}

		// Check if token is valid
		if !token.Valid {
			return jsend.SendUnauthorized(c, map[string]string{
				"message": "Invalid token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return jsend.SendUnauthorized(c, map[string]string{
				"message": "Invalid token claims",
			})
		}

		// Set user context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	})
}

// Optional middleware - only validates token if present
func NewOptionalAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Next()
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Next()
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			return c.Next()
		}

		claims, ok := token.Claims.(*Claims)
		if ok {
			c.Locals("userID", claims.UserID)
			c.Locals("email", claims.Email)
		}

		return c.Next()
	}
}
