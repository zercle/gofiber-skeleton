package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

func TestRecovery(t *testing.T) {
	t.Run("recovers from panic", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.Recovery())
		app.Get("/panic", func(c *fiber.Ctx) error {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("does not interfere with normal requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.Recovery())
		app.Get("/normal", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/normal", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("recovers from nil pointer dereference", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.Recovery())
		app.Get("/nil-pointer", func(c *fiber.Ctx) error {
			var ptr *string
			_ = *ptr // This will panic
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/nil-pointer", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("recovers from index out of range", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.Recovery())
		app.Get("/out-of-range", func(c *fiber.Ctx) error {
			slice := []int{1, 2, 3}
			_ = slice[10] // This will panic
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/out-of-range", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("continues serving after panic", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.Recovery())
		app.Get("/panic", func(c *fiber.Ctx) error {
			panic("test panic")
		})
		app.Get("/normal", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// First request panics
		req1 := httptest.NewRequest("GET", "/panic", nil)
		resp1, err := app.Test(req1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp1.StatusCode)

		// Second request should work normally
		req2 := httptest.NewRequest("GET", "/normal", nil)
		resp2, err := app.Test(req2)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
	})
}
