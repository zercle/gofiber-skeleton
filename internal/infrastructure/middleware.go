package infrastructure

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware creates JWT authentication middleware
func JWTMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip authentication for public endpoints
		if isPublicEndpoint(c.Path()) {
			return c.Next()
		}

		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Check if token is valid
		if !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Set user information in context
		c.Locals("user_id", claims["user_id"])
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

// AdminMiddleware creates admin-only authorization middleware
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}
		return c.Next()
	}
}

// isPublicEndpoint checks if the endpoint is public (no authentication required)
func isPublicEndpoint(path string) bool {
	publicEndpoints := []string{
		"/api/v1/register",
		"/api/v1/login",
		"/api/v1/products", // GET /products is public
	}

	for _, endpoint := range publicEndpoints {
		if path == endpoint {
			return true
		}
	}

	// GET /products/{id} is also public
	if strings.HasPrefix(path, "/api/v1/products/") && !strings.Contains(path, "/") {
		return true
	}

	return false
}