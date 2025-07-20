package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"golang.org/x/crypto/bcrypt"

	"gofiber-skeleton/internal/entities"
	"gofiber-skeleton/mocks"
)

func TestUserUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	username := "testuser"
	password := "testpassword"

	u7, err := uuid.NewV7()
	assert.NoError(t, err)
	expectedUser := &entities.User{
		ID:        u7.String(),
		Username:  username,
		Password:  "hashedpassword", // This will be hashed by the usecase
		CreatedAt: time.Now(),
	}

	mockRepo.EXPECT().CreateUser(gomock.Any(), username, gomock.Any()).Return(expectedUser, nil)

	var user *entities.User
	user, err = usecase.Register(context.Background(), username, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.NotEqual(t, password, user.Password) // Password should be hashed
}

func TestUserUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	username := "testuser"
	password := "testpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	u7, err := uuid.NewV7()
	assert.NoError(t, err)
	existingUser := &entities.User{
		ID:        u7.String(),
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Test successful login
	mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(existingUser, nil)
	var user *entities.User
	user, err = usecase.Login(context.Background(), username, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, existingUser.ID, user.ID)

	// Test invalid password
	mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(existingUser, nil)
	user, err = usecase.Login(context.Background(), username, "wrongpassword")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Test user not found
	mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, errors.New("user not found"))
	user, err = usecase.Login(context.Background(), username, password)
	assert.Error(t, err)
	assert.Nil(t, user)
}
