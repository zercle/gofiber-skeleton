package users_test

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zercle/gofiber-skelton/internal/users"
	"github.com/zercle/gofiber-skelton/mocks"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

func TestGetUserUsecase(t *testing.T) {
	var mockUser models.User
	gofakeit.Struct(&mockUser)

	mockUserRepo := new(mocks.UserRepository)
	mockUserUsecase := new(mocks.UserUsecase)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetUser", mock.AnythingOfType("string")).Return(mockUser, nil).Once()

		usecase := users.InitUserUsecase(mockUserRepo)

		result, err := usecase.GetUser(mockUser.Id)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, result)

		mockUserUsecase.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("GetUser", mock.AnythingOfType("string")).Return(models.User{}, errors.New("call error")).Once()

		usecase := users.InitUserUsecase(mockUserRepo)

		result, err := usecase.GetUser(mockUser.Id)

		assert.Error(t, err)
		assert.Equal(t, models.User{}, result)

		mockUserUsecase.AssertExpectations(t)
	})

}

func TestGetUsersUsecase(t *testing.T) {
	var mockUser models.User
	gofakeit.Struct(&mockUser)

	mockUsers := []models.User{mockUser}

	mockUserRepo := new(mocks.UserRepository)
	mockUserUsecase := new(mocks.UserUsecase)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetUsers", mock.Anything).Return(mockUsers, nil).Once()

		usecase := users.InitUserUsecase(mockUserRepo)

		result, err := usecase.GetUsers(mockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUsers, result)

		mockUserUsecase.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("GetUsers", mock.Anything).Return([]models.User{}, errors.New("call error")).Once()

		usecase := users.InitUserUsecase(mockUserRepo)

		result, err := usecase.GetUsers(mockUser)

		assert.Error(t, err)
		assert.Empty(t, result)

		mockUserUsecase.AssertExpectations(t)
	})

}
