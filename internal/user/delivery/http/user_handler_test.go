package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	userHandler "gofiber-skeleton/internal/user/delivery/http"
	mock_usecase "gofiber-skeleton/internal/user/usecase/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	handler := userHandler.NewHTTPUserHandler(mockUserUseCase)

	app := fiber.New()
	app.Post("/api/users/register", handler.Register)

	// Test case 1: Successful registration
	username := "testuser"
	password := "testpassword"
	registerPayload := map[string]string{
		"username": username,
		"password": password,
	}

	mockUserUseCase.EXPECT().Register(gomock.Any(), username, password).Return(nil, nil).Times(1)

	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]string
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User registered successfully", response["message"])

	// Test case 2: Invalid request body
	req = httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test case 3: Use case error
	useCaseError := errors.New("registration failed")
	mockUserUseCase.EXPECT().Register(gomock.Any(), username, password).Return(nil, useCaseError).Times(1)

	body, _ = json.Marshal(registerPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	handler := userHandler.NewHTTPUserHandler(mockUserUseCase)

	app := fiber.New()
	app.Post("/api/users/login", handler.Login)

	// Test case 1: Successful login
	var err error // Declare err here
	username := "testuser"
	password := "testpassword"
	token := "mock_jwt_token"
	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}

	mockUserUseCase.EXPECT().Login(gomock.Any(), username, password).Return(token, nil).Times(1)

	body, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, token, data["token"])

	// Test case 2: Invalid request body
	req = httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test case 3: Use case error (e.g., invalid credentials)
	useCaseError := errors.New("invalid credentials")
	mockUserUseCase.EXPECT().Login(gomock.Any(), username, password).Return("", useCaseError).Times(1)

	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "Invalid credentials", response["message"])
}
