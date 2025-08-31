package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/user"
	"github.com/zercle/gofiber-skeleton/internal/user/mock"
	"go.uber.org/mock/gomock"
)

func TestUserUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo)

	t.Run("successful registration", func(t *testing.T) {
		payload := user.RegisterPayload{
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedUser := user.User{
			ID:        uuid.New(),
			Email:     payload.Email,
			Password:  "$2a$10$hashedpassword", // Mock hashed password
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.EXPECT().
			CreateUser(gomock.Any()).
			DoAndReturn(func(u user.User) (user.User, error) {
				// Verify that the password was hashed
				assert.NotEqual(t, payload.Password, u.Password)
				assert.NotEmpty(t, u.Password)
				// Return expected user with ID set
				return expectedUser, nil
			}).
			Times(1)

		result, err := userUsecase.Register(payload)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.Password, result.Password)
	})

	t.Run("repository error", func(t *testing.T) {
		payload := user.RegisterPayload{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.EXPECT().
			CreateUser(gomock.Any()).
			Return(user.User{}, errors.New("database error")).
			Times(1)

		result, err := userUsecase.Register(payload)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, result)
		assert.Contains(t, err.Error(), "database error")
	})
}