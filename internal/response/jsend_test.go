package response_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/internal/response"
)

func setupTestApp() *fiber.App {
	return fiber.New()
}

func parseResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	return result
}

func TestSuccess(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, fiber.Map{
			"message": "success",
			"count":   42,
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	result := parseResponse(t, resp)
	assert.Equal(t, "success", result["status"])

	data, ok := result["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "success", data["message"])
	assert.Equal(t, float64(42), data["count"])
}

func TestFail(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Fail(c, fiber.StatusBadRequest, fiber.Map{
			"email": "invalid email format",
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	result := parseResponse(t, resp)
	assert.Equal(t, "fail", result["status"])

	data, ok := result["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "invalid email format", data["email"])
}

func TestError(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusInternalServerError, "database connection failed", 5001)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	result := parseResponse(t, resp)
	assert.Equal(t, "error", result["status"])
	assert.Equal(t, "database connection failed", result["message"])
	assert.Equal(t, float64(5001), result["code"])
}

func TestSuccessWithNilData(t *testing.T) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, nil)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	result := parseResponse(t, resp)
	assert.Equal(t, "success", result["status"])
	assert.Nil(t, result["data"])
}
