package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/mocks"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
)

func TestUserUsecase_Register(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	tests := []struct {
		name    string
		request *entity.CreateUserRequest
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "successful registration",
			request: &entity.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			setup: func() {
				mockRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			request: &entity.CreateUserRequest{
				Email:    "existing@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			setup: func() {
				mockRepo.On("EmailExists", mock.Anything, "existing@example.com").Return(true, nil)
			},
			wantErr: true,
			errType: entity.ErrEmailExists,
		},
		{
			name: "invalid email",
			request: &entity.CreateUserRequest{
				Email:    "invalid-email",
				Password: "password123",
				FullName: "Test User",
			},
			setup:   func() {},
			wantErr: true,
			errType: entity.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockRepo.ExpectedCalls = nil

			// Setup test
			tt.setup()

			// Execute
			result, err := userUsecase.Register(context.Background(), tt.request)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Email, result.Email)
				assert.Equal(t, tt.request.FullName, result.FullName)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Create a test user
	testUser := &entity.DomainUser{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=1,p=4$c29tZXNhbHQ$testhash",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name    string
		request *entity.LoginRequest
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name: "successful login",
			request: &entity.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			request: &entity.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, entity.ErrUserNotFound)
			},
			wantErr: true,
			errType: entity.ErrInvalidPassword,
		},
		{
			name: "invalid email format",
			request: &entity.LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			setup:   func() {},
			wantErr: true,
			errType: entity.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockRepo.ExpectedCalls = nil

			// Setup test
			tt.setup()

			// Execute
			result, err := userUsecase.Login(context.Background(), tt.request)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Token)
				assert.Equal(t, testUser.Email, result.User.Email)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_GetProfile(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: "24h",
		},
	}

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)
	do.ProvideValue(injector, cfg)

	userUsecase, err := usecase.NewUserUsecase(injector)
	require.NoError(t, err)

	// Create a test user
	testUser := &entity.DomainUser{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=1,p=4$c29tZXNhbHQ$testhash",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name    string
		userID  uuid.UUID
		setup   func()
		wantErr bool
		errType error
	}{
		{
			name:   "successful profile retrieval",
			userID: testUser.ID,
			setup: func() {
				mockRepo.On("GetByID", mock.Anything, testUser.ID).Return(testUser, nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			setup: func() {
				mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, entity.ErrUserNotFound)
			},
			wantErr: true,
			errType: entity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockRepo.ExpectedCalls = nil

			// Setup test
			tt.setup()

			// Execute
			result, err := userUsecase.GetProfile(context.Background(), tt.userID)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testUser.Email, result.Email)
				assert.Equal(t, testUser.FullName, result.FullName)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
