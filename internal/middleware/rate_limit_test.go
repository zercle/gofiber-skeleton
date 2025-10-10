package middleware_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

func TestRateLimit(t *testing.T) {
	t.Run("allows requests within limit", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RateLimit(5, 1*time.Minute))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Send 5 requests (within limit)
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}
	})

	t.Run("blocks requests exceeding limit", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RateLimit(3, 1*time.Minute))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Send 3 requests (at limit)
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		// 4th request should be blocked
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)

		// Verify error message
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response["status"])
		assert.Contains(t, response["message"], "Rate limit exceeded")
	})

	t.Run("rate limit configuration works", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RateLimit(2, 1*time.Minute))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// First 2 requests should succeed
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		// 3rd request should be blocked
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})
}

func TestAuthRateLimit(t *testing.T) {
	t.Run("auth rate limit configuration", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.AuthRateLimit())
		app.Post("/auth/login", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Send 5 requests (at limit for auth: 5 requests per 15 minutes)
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("POST", "/auth/login", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		// 6th request should be blocked
		req := httptest.NewRequest("POST", "/auth/login", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})
}

func TestAPIRateLimit(t *testing.T) {
	t.Run("API rate limit configuration", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.APIRateLimit())
		app.Get("/api/data", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Send 100 requests (at limit for API: 100 requests per minute)
		for i := 0; i < 100; i++ {
			req := httptest.NewRequest("GET", "/api/data", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		// 101st request should be blocked
		req := httptest.NewRequest("GET", "/api/data", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})
}

func TestRateLimitHeaders(t *testing.T) {
	t.Run("rate limit headers present", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RateLimit(10, 1*time.Minute))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Check if rate limit headers are present (Fiber limiter adds these)
		// X-RateLimit-* headers may be present depending on Fiber version
		// This test ensures the middleware is working
		assert.NotEmpty(t, resp.Header.Get("X-RateLimit-Limit") + resp.Header.Get("RateLimit-Limit"))
	})
}
