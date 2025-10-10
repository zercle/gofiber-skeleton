package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

func TestSecurity(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Security())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tests := []struct {
		name          string
		headerName    string
		expectedValue string
		checkContains bool
	}{
		{
			name:          "XSS Protection header",
			headerName:    "X-Xss-Protection",
			expectedValue: "1; mode=block",
			checkContains: false,
		},
		{
			name:          "Content Type Nosniff header",
			headerName:    "X-Content-Type-Options",
			expectedValue: "nosniff",
			checkContains: false,
		},
		{
			name:          "X-Frame-Options header",
			headerName:    "X-Frame-Options",
			expectedValue: "SAMEORIGIN",
			checkContains: false,
		},
		{
			name:          "Referrer-Policy header",
			headerName:    "Referrer-Policy",
			expectedValue: "no-referrer",
			checkContains: false,
		},
		{
			name:          "Content-Security-Policy header",
			headerName:    "Content-Security-Policy",
			expectedValue: "default-src 'self'",
			checkContains: false,
		},
		{
			name:          "Cross-Origin-Embedder-Policy header",
			headerName:    "Cross-Origin-Embedder-Policy",
			expectedValue: "require-corp",
			checkContains: false,
		},
		{
			name:          "Cross-Origin-Opener-Policy header",
			headerName:    "Cross-Origin-Opener-Policy",
			expectedValue: "same-origin",
			checkContains: false,
		},
		{
			name:          "Cross-Origin-Resource-Policy header",
			headerName:    "Cross-Origin-Resource-Policy",
			expectedValue: "same-origin",
			checkContains: false,
		},
		{
			name:          "X-DNS-Prefetch-Control header",
			headerName:    "X-Dns-Prefetch-Control",
			expectedValue: "off",
			checkContains: false,
		},
		{
			name:          "X-Download-Options header",
			headerName:    "X-Download-Options",
			expectedValue: "noopen",
			checkContains: false,
		},
		{
			name:          "X-Permitted-Cross-Domain-Policies header",
			headerName:    "X-Permitted-Cross-Domain-Policies",
			expectedValue: "none",
			checkContains: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			headerValue := resp.Header.Get(tt.headerName)
			if tt.checkContains {
				assert.Contains(t, headerValue, tt.expectedValue)
			} else {
				assert.Equal(t, tt.expectedValue, headerValue)
			}
		})
	}

	t.Run("critical security headers present", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Verify key security headers are set
		securityHeaders := []string{
			"X-Xss-Protection",
			"X-Content-Type-Options",
			"X-Frame-Options",
			"Strict-Transport-Security",
			"Content-Security-Policy",
			"Referrer-Policy",
			"Cross-Origin-Embedder-Policy",
			"Cross-Origin-Opener-Policy",
			"Cross-Origin-Resource-Policy",
			"X-Dns-Prefetch-Control",
		}

		presentHeaders := 0
		for _, header := range securityHeaders {
			if resp.Header.Get(header) != "" {
				presentHeaders++
			}
		}

		assert.GreaterOrEqual(t, presentHeaders, 8, "Expected at least 8 key security headers to be set")
	})
}
