package users_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/mocks"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/users"
)

func TestGetUserHandler(t *testing.T) {
	var mockUser models.User

	gofakeit.Struct(&mockUser)

	mockResponse := helpers.ResponseForm{
		Success: true,
		Result: map[string]interface{}{
			"user": mockUser,
		},
	}

	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetUser", mockUser.Id).Return(mockUser, nil)

	app := fiber.New()
	req, err := http.NewRequest("GET", "/api/v1/user/"+mockUser.Id, nil)
	assert.NoError(t, err)

	mockUCase := users.NewUserUsecase(mockRepo)
	users.NewUserHandler(app.Group("/api/v1/user"), mockUCase)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockJson, err := json.Marshal(mockResponse)
	assert.NoError(t, err)

	var bodyBytes bytes.Buffer
	io.Copy(&bodyBytes, resp.Body)

	assert.JSONEq(t, string(mockJson), bodyBytes.String())

	mockRepo.AssertExpectations(t)
}
