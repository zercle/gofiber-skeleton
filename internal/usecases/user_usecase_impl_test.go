package usecases

import (
	"context"
	"testing"

	"gofiber-skeleton/internal/entities"
	"gofiber-skeleton/internal/usecases/mocks"

	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUseCase(mockUserRepo, "test-secret", 1)

	username := "testuser"
	password := "password"

	mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)

	user, err := userUseCase.Register(context.Background(), username, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	assert.NoError(t, err)
}

func TestUserUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUseCase(mockUserRepo, "test-secret", 1)

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(&entities.User{Username: username, Password: string(hashedPassword)}, nil)

	token, err := userUseCase.Login(context.Background(), username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
