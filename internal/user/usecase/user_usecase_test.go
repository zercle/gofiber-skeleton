package usecase_test

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/user"
	mock_repository "gofiber-skeleton/internal/user/repository/mocks"
	"gofiber-skeleton/internal/user/usecase"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	duration, _ := time.ParseDuration("1h")
	userUseCase := usecase.NewUserUseCase(mockUserRepo, "test_secret", duration)

	ctx := context.Background()
	username := "testuser"
	password := "testpassword"

	// Test case 1: Successful user registration
	mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil).Times(1)

	registeredUser, err := userUseCase.Register(ctx, username, password)
	assert.NoError(t, err)
	assert.NotNil(t, registeredUser)
	assert.Equal(t, username, registeredUser.Username)
	assert.NotEmpty(t, registeredUser.Password) // Password should be hashed

	// Test case 2: Error creating user in repository
	repoError := errors.New("failed to create user in repository")
	mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(repoError).Times(1)

	registeredUser, err = userUseCase.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Nil(t, registeredUser)
	assert.Equal(t, repoError, err)
}

func TestUserUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	jwtSecret := "test_secret"
	jwtExpirationStr := "1h"
	jwtExpiration, _ := time.ParseDuration(jwtExpirationStr)
	userUseCase := usecase.NewUserUseCase(mockUserRepo, jwtSecret, jwtExpiration)

	ctx := context.Background()
	username := "testuser"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	testUser := &user.User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
	}

	// Capture the current time before generating the token

	// Test case 1: Successful login
	mockUserRepo.EXPECT().GetUserByUsername(ctx, username).Return(testUser, nil).Times(1)

	now := time.Now()
	token, err := userUseCase.Login(ctx, username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, testUser.ID.String(), claims["sub"])
	// Convert claims["exp"] to time.Time for more accurate comparison
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	assert.WithinDuration(t, now.Add(jwtExpiration), expirationTime, 2*time.Second) // Allow a 2-second margin

	// Test case 2: Invalid password
	mockUserRepo.EXPECT().GetUserByUsername(ctx, username).Return(testUser, nil).Times(1)

	token, err = userUseCase.Login(ctx, username, "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "crypto/bcrypt: hashedPassword is not the hash of the given password", err.Error())

	// Test case 3: User not found (repository error)
	repoError := errors.New("user not found")
	mockUserRepo.EXPECT().GetUserByUsername(ctx, "nonexistent").Return(nil, repoError).Times(1)

	token, err = userUseCase.Login(ctx, "nonexistent", password)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, repoError, err)
}
