package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/mocks"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
)


func TestUserUsecase_Register_EmailExists(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	req := &entity.CreateUserRequest{
		Email:    "existing@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	// Setup expectations
	mockRepo.EXPECT().EmailExists(gomock.Any(), "existing@example.com").Return(true, nil)

	// Execute
	result, err := userUsecase.Register(context.Background(), req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrEmailExists, err)
}

func TestUserUsecase_Login_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Create test user with hashed password
	user := &entity.DomainUser{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Test data
	req := &entity.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Setup expectations
	mockRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(user, nil)

	// Execute
	result, err := userUsecase.Login(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Token)
	assert.NotNil(t, result.User)
	assert.Equal(t, user.Email, result.User.Email)
	assert.Equal(t, user.FullName, result.User.FullName)
}

func TestUserUsecase_Login_UserNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	req := &entity.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Setup expectations
	mockRepo.EXPECT().GetByEmail(gomock.Any(), "nonexistent@example.com").Return(nil, entity.ErrUserNotFound)

	// Execute
	result, err := userUsecase.Login(context.Background(), req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrInvalidPassword, err)
}

func TestUserUsecase_GetProfile_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	user := &entity.DomainUser{
		ID:        userID,
		Email:     "test@example.com",
		FullName:  "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)

	// Execute
	result, err := userUsecase.GetProfile(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.FullName, result.FullName)
	assert.Equal(t, user.IsActive, result.IsActive)
}

func TestUserUsecase_GetProfile_UserNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(nil, entity.ErrUserNotFound)

	// Execute
	result, err := userUsecase.GetProfile(context.Background(), userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrUserNotFound, err)
}

func TestUserUsecase_UpdateProfile_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	existingUser := &entity.DomainUser{
		ID:        userID,
		Email:     "test@example.com",
		FullName:  "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &entity.UpdateUserRequest{
		Email:    "updated@example.com",
		FullName: "Updated User",
	}

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(existingUser, nil)
	mockRepo.EXPECT().EmailExists(gomock.Any(), "updated@example.com").Return(false, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	result, err := userUsecase.UpdateProfile(context.Background(), userID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, req.FullName, result.FullName)
}

func TestUserUsecase_UpdateProfile_EmailExists(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	existingUser := &entity.DomainUser{
		ID:        userID,
		Email:     "test@example.com",
		FullName:  "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &entity.UpdateUserRequest{
		Email: "existing@example.com",
	}

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(existingUser, nil)
	mockRepo.EXPECT().EmailExists(gomock.Any(), "existing@example.com").Return(true, nil)

	// Execute
	result, err := userUsecase.UpdateProfile(context.Background(), userID, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrEmailExists, err)
}

func TestUserUsecase_ListUsers(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	users := []*entity.DomainUser{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			FullName:  "User One",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     "user2@example.com",
			FullName:  "User Two",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Setup expectations
	mockRepo.EXPECT().List(gomock.Any(), 10, 0).Return(users, nil)

	// Execute
	result, err := userUsecase.ListUsers(context.Background(), 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].Email, result[0].Email)
	assert.Equal(t, users[1].Email, result[1].Email)
}

func TestUserUsecase_DeactivateUser_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()

	// Setup expectations
	mockRepo.EXPECT().Deactivate(gomock.Any(), userID).Return(nil)

	// Execute
	err := userUsecase.DeactivateUser(context.Background(), userID)

	// Assert
	require.NoError(t, err)
}

func TestUserUsecase_DeactivateUser_Error(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		AppName: "test-app",
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()

	// Setup expectations
	mockRepo.EXPECT().Deactivate(gomock.Any(), userID).Return(entity.ErrUserNotFound)

	// Execute
	err := userUsecase.DeactivateUser(context.Background(), userID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, entity.ErrUserNotFound, err)
}
