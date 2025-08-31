package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/user"
	"github.com/zercle/gofiber-skeleton/internal/user/mock"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock.NewMockUserUsecase(ctrl)
	handler := NewUserHandler(mockUsecase)

	app := fiber.New()
	app.Post("/register", handler.Register)

	t.Run("successful registration", func(t *testing.T) {
		payload := user.RegisterPayload{
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedUser := user.User{
			ID:    uuid.New(),
			Email: payload.Email,
		}

		mockUsecase.EXPECT().
			Register(payload).
			Return(expectedUser, nil).
			Times(1)

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("invalid json body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		payload := user.RegisterPayload{
			Email:    "invalid-email",
			Password: "123", // Too short
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}