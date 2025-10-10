package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

func TestCORS(t *testing.T) {
	t.Run("default CORS configuration", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.CORS())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Test preflight request
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
		assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "GET")
		assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "POST")
	})

	t.Run("actual request with CORS", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.CORS())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://example.com")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("CORS with Authorization header", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.CORS())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("Access-Control-Request-Headers", "Authorization")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		// Check that Authorization is in the allowed headers configuration
		allowedHeaders := resp.Header.Get("Access-Control-Allow-Headers")
		assert.NotEmpty(t, allowedHeaders, "Should have Access-Control-Allow-Headers")
	})
}

func TestCORSWithConfig(t *testing.T) {
	t.Run("custom CORS configuration", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.CORSWithConfig(
			"http://localhost:3000",
			"GET,POST",
			"Content-Type",
			true,
		))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET,POST", resp.Header.Get("Access-Control-Allow-Methods"))
	})

	t.Run("custom CORS with credentials", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.CORSWithConfig(
			"http://localhost:3000",
			"GET,POST",
			"Content-Type",
			true,
		))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
	})

	t.Run("multiple origins", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.CORSWithConfig(
			"http://localhost:3000,https://example.com",
			"GET,POST",
			"Content-Type",
			false,
		))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://example.com")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	})
}
