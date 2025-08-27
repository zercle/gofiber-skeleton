package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/domain/mock"
	userhandler "github.com/zercle/gofiber-skeleton/internal/user/handler"
)

func setupUserIntegrationTest(t *testing.T) (*fiber.App, *mock.MockUserUseCase) {
	ctrl := gomock.NewController(t)
	// defer ctrl.Finish() // Don't defer here, let individual tests manage it if needed.

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	userHandler := userhandler.NewUserHandler(mockUserUseCase)

	app := fiber.New()
	app.Post("/api/v1/register", userHandler.Register)
	app.Post("/api/v1/login", userHandler.Login)

	return app, mockUserUseCase
}

func TestUserIntegration_Register(t *testing.T) {
	app, mockUserUseCase := setupUserIntegrationTest(t)

	registerInput := userhandler.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	mockTime := time.Now()

	t.Run("successful user registration", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:           uuid.New().String(),
			Username:     registerInput.Username,
			PasswordHash: "any_hashed_password",
			Role:         domain.RoleCustomer,
			CreatedAt:    mockTime,
			UpdatedAt:    mockTime,
		}

		mockUserUseCase.EXPECT().Register(
			registerInput.Username,
			registerInput.Password,
			"",
		).Return(expectedUser, nil)

		body, _ := json.Marshal(registerInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		data := responseBody["data"].(map[string]any)
		assert.Equal(t, "User registered successfully", data["message"])
		user := data["user"].(map[string]any)
		assert.Equal(t, registerInput.Username, user["username"])
	})

	t.Run("registration with existing username", func(t *testing.T) {
		mockUserUseCase.EXPECT().Register(
			registerInput.Username,
			registerInput.Password,
			"",
		).Return(nil, errors.New("username already exists"))

		body, _ := json.Marshal(registerInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "username already exists", responseBody["message"])
	})
}

func TestUserIntegration_Login(t *testing.T) {
	app, mockUserUseCase := setupUserIntegrationTest(t)

	loginInput := userhandler.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	t.Run("successful user login", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:           uuid.New().String(),
			Username:     loginInput.Username,
			PasswordHash: "hashed_password_mock",
			Role:         domain.RoleCustomer,
		}
		expectedToken := "mock-jwt-token"

		mockUserUseCase.EXPECT().Login(
			loginInput.Username,
			loginInput.Password,
		).Return(expectedToken, expectedUser, nil)

		body, _ := json.Marshal(loginInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		assert.Contains(t, responseBody, "data")
		data := responseBody["data"].(map[string]any)
		assert.Equal(t, "Login successful", data["message"])
		assert.Equal(t, expectedToken, data["token"])
		user := data["user"].(map[string]any)
		assert.Equal(t, loginInput.Username, user["username"])
	})

	t.Run("login with invalid credentials", func(t *testing.T) {
		invalidLoginInput := userhandler.LoginRequest{
			Username: "testuser",
			Password: "wrong_password",
		}

		mockUserUseCase.EXPECT().Login(
			invalidLoginInput.Username,
			invalidLoginInput.Password,
		).Return("", nil, errors.New("invalid credentials"))

		body, _ := json.Marshal(invalidLoginInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
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
