package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gofiber-skeleton/internal/entities"
	deliveryhttp "gofiber-skeleton/internal/delivery/http"
	"gofiber-skeleton/internal/usecases/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_ = context.Background() // Dummy usage to prevent "context imported and not used" error

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	handler := deliveryhttp.NewHTTPUserHandler(mockUserUseCase)

	app := fiber.New()
	app.Post("/register", handler.Register)

	tests := []struct {
		name           string
		requestBody    map[string]string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful registration",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().Register(gomock.Any(), "testuser", "password123").Return(&entities.User{
					ID:        [16]byte{},
					Username:  "testuser",
					Password:  "hashedpassword",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":"success","message":"User registered successfully"}`,
		},
		{
			name: "registration with existing username",
			requestBody: map[string]string{
				"username": "existinguser",
				"password": "password123",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().Register(gomock.Any(), "existinguser", "password123").Return(nil, errors.New("username already exists"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","message":"Username already exists"}`,
		},
		{
			name: "invalid request body",
			requestBody: map[string]string{
				"username": "testuser",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","message":"Username and password cannot be empty"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBody, string(respBody))
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	handler := deliveryhttp.NewHTTPUserHandler(mockUserUseCase)

	app := fiber.New()
	app.Post("/login", handler.Login)

	tests := []struct {
		name           string
		requestBody    map[string]string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful login",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().Login(gomock.Any(), "testuser", "password123").Return("mock_jwt_token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success","data":{"token":"mock_jwt_token"}}`,
		},
		{
			name: "login with invalid credentials",
			requestBody: map[string]string{
				"username": "wronguser",
				"password": "wrongpassword",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().Login(gomock.Any(), "wronguser", "wrongpassword").Return("", errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":"error","message":"Invalid credentials"}`,
		},
		{
			name: "invalid request body",
			requestBody: map[string]string{
				"username": "testuser",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","message":"Username and password cannot be empty"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBody, string(respBody))
		})
	}
}