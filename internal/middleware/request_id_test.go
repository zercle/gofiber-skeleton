package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

func TestRequestID(t *testing.T) {
	t.Run("generates request ID when not provided", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RequestID())
		app.Get("/test", func(c *fiber.Ctx) error {
			requestID := c.Locals("requestid").(string)
			return c.SendString(requestID)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		requestID := resp.Header.Get("X-Request-ID")
		assert.NotEmpty(t, requestID)

		// Verify it's a valid UUID
		_, err = uuid.Parse(requestID)
		assert.NoError(t, err, "Request ID should be a valid UUID")
	})

	t.Run("uses provided request ID from header", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RequestID())
		app.Get("/test", func(c *fiber.Ctx) error {
			requestID := c.Locals("requestid").(string)
			return c.SendString(requestID)
		})

		providedID := "custom-request-id-12345"
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", providedID)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, providedID, resp.Header.Get("X-Request-ID"))
	})

	t.Run("request ID available in context locals", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RequestID())
		app.Get("/test", func(c *fiber.Ctx) error {
			requestID := c.Locals("requestid")
			assert.NotNil(t, requestID)
			assert.IsType(t, "", requestID)
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("different requests get different IDs", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RequestID())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// First request
		req1 := httptest.NewRequest("GET", "/test", nil)
		resp1, err := app.Test(req1)
		assert.NoError(t, err)
		requestID1 := resp1.Header.Get("X-Request-ID")

		// Second request
		req2 := httptest.NewRequest("GET", "/test", nil)
		resp2, err := app.Test(req2)
		assert.NoError(t, err)
		requestID2 := resp2.Header.Get("X-Request-ID")

		// IDs should be different
		assert.NotEqual(t, requestID1, requestID2)
	})

	t.Run("request ID persists through handler chain", func(t *testing.T) {
		var capturedID string

		app := fiber.New()
		app.Use(middleware.RequestID())
		app.Use(func(c *fiber.Ctx) error {
			capturedID = c.Locals("requestid").(string)
			return c.Next()
		})
		app.Get("/test", func(c *fiber.Ctx) error {
			handlerID := c.Locals("requestid").(string)
			assert.Equal(t, capturedID, handlerID, "Request ID should be the same across handler chain")
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, capturedID)
	})
}
