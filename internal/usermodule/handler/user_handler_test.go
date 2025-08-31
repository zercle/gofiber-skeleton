package userhandler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/internal/usermodule"
	usermock "github.com/zercle/gofiber-skeleton/internal/usermodule/mock"
	"go.uber.org/mock/gomock"
)

// mockAuthMiddleware is a Fiber middleware that sets a user ID in c.Locals for testing.
/*
func mockAuthMiddleware(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	}
}
*/

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := usermock.NewMockUserUseCase(ctrl)
	app := fiber.New()
	handler := NewUserHandler(mockUserUseCase)
	app.Post("/api/v1/register", handler.Register)

	registerReq := RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Role:     "customer",
	}
	bodyBytes, _ := json.Marshal(registerReq)

	t.Run("successful registration", func(t *testing.T) {
		expectedUser := &usermodule.User{
			ID:        uuid.New().String(),
			Username:  registerReq.Username,
			Role:      registerReq.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockUserUseCase.EXPECT().Register(registerReq.Username, registerReq.Password, registerReq.Role).Return(expectedUser, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		data := responseBody["data"].(map[string]any)
		assert.Equal(t, "User registered successfully", data["message"])
		assert.NotNil(t, data["user"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader([]byte(`{"username":123}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request body", responseBody["message"])
	})

	t.Run("missing username or password", func(t *testing.T) {
		invalidReq := RegisterRequest{Username: "testuser", Password: ""}
		invalidBodyBytes, _ := json.Marshal(invalidReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(invalidBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "Username and password are required", responseBody["message"])
	})

	t.Run("usecase returns error", func(t *testing.T) {
		mockUserUseCase.EXPECT().Register(registerReq.Username, registerReq.Password, registerReq.Role).Return(nil, errors.New("username already exists"))

		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "username already exists", responseBody["message"])
	})
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := usermock.NewMockUserUseCase(ctrl)
	app := fiber.New()
	handler := NewUserHandler(mockUserUseCase)
	app.Post("/api/v1/login", handler.Login)

	loginReq := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(loginReq)

	t.Run("successful login", func(t *testing.T) {
		expectedToken := "mock_jwt_token"
		expectedUser := &usermodule.User{
			ID:       uuid.New().String(),
			Username: loginReq.Username,
			Role:     usermodule.RoleCustomer,
		}
		mockUserUseCase.EXPECT().Login(loginReq.Username, loginReq.Password).Return(expectedToken, expectedUser, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		data := responseBody["data"].(map[string]any)
		assert.Equal(t, "Login successful", data["message"])
		assert.Equal(t, expectedToken, data["token"])
		user := data["user"].(map[string]any)
		assert.Equal(t, expectedUser.Username, user["username"])
		assert.Equal(t, expectedUser.ID, user["id"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{"username":123}`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "Invalid request body", responseBody["message"])
	})

	t.Run("missing username or password", func(t *testing.T) {
		invalidReq := LoginRequest{Username: "testuser", Password: ""}
		invalidBodyBytes, _ := json.Marshal(invalidReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(invalidBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "Username and password are required", responseBody["message"])
	})

	t.Run("usecase returns error", func(t *testing.T) {
		mockUserUseCase.EXPECT().Login(loginReq.Username, loginReq.Password).Return("", (*usermodule.User)(nil), errors.New("invalid credentials"))

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "invalid credentials", responseBody["message"])
	})
}
