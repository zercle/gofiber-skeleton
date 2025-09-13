package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock" // Correct import for gomock
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/entities"
	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/mocks" // Import the generated mock
	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/models"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)


// TestAuthUseCase_Register tests the Register method of AuthUseCase
func TestAuthUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t) // Initialize gomock controller
	defer ctrl.Finish()             // Ensure controller is cleaned up

	mockRepo := mocks.NewMockUserRepository(ctrl) // Use generated mock
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "testsecret",
			ExpiresIn: "1h",
		},
	}
	uc := NewAuthUseCase(mockRepo, cfg)

	req := &models.RegisterRequest{
		Email:           "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		FirstName:       "John",
		LastName:        "Doe",
	}

	// Test case 1: Successful registration
	mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(nil, nil).Times(1)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entities.User{
		ID:              uuid.New(),
		Email:           req.Email,
		PasswordHash:    "hashedpassword", // This will be set by the usecase
		FirstName:       "John",
		LastName:        "Doe",
		IsActive:        true,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil).Times(1)

	resp, err := uc.Register(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, req.Email, resp.User.Email)

	// Test case 2: Email already registered
	mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(&entities.User{}, nil).Times(1)

	resp, err = uc.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "email already registered")

	// Test case 3: Repository create error
	mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(nil, nil).Times(1)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)

	resp, err = uc.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "failed to create user")
}

// TestAuthUseCase_Login tests the Login method of AuthUseCase
func TestAuthUseCase_Login(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "testsecret",
			ExpiresIn: "1h",
		},
	}

	t.Run("Successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &entities.User{
			ID:              uuid.New(),
			Email:           req.Email,
			PasswordHash:    string(hashedPassword),
			FirstName:       "John",
			LastName:        "Doe",
			IsActive:        true,
			IsEmailVerified: true,
			LastLoginAt:     nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(user, nil).Times(1)
		mockRepo.EXPECT().UpdateLastLogin(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		resp, err := uc.Login(context.Background(), req)
		assert.NoError(t, err)
		if resp != nil {
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.Token)
			assert.Equal(t, req.Email, resp.User.Email)
		} else {
			assert.NotNil(t, err)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		req := &models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(nil, errors.New("user not found")).Times(1)

		resp, err := uc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "invalid credentials")
	})

	t.Run("Invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		user := &entities.User{
			ID:              uuid.New(),
			Email:           req.Email,
			PasswordHash:    hashPassword("password123"), // Correct hashed password
			FirstName:       "John",
			LastName:        "Doe",
			IsActive:        true,
			IsEmailVerified: true,
			LastLoginAt:     nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(user, nil).Times(1)

		resp, err := uc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "invalid credentials")
	})

	t.Run("Account deactivated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		user := &entities.User{
			ID:              uuid.New(),
			Email:           req.Email,
			PasswordHash:    hashPassword("password123"),
			IsActive:        false, // Account is deactivated
			FirstName:       "John",
			LastName:        "Doe",
			IsEmailVerified: true,
			LastLoginAt:     nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(user, nil).Times(1)

		resp, err := uc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "account is deactivated")
	})

	t.Run("Failed to update last login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		user := &entities.User{
			ID:              uuid.New(),
			Email:           req.Email,
			PasswordHash:    hashPassword("password123"),
			IsActive:        true,
			FirstName:       "John",
			LastName:        "Doe",
			IsEmailVerified: true,
			LastLoginAt:     nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(user, nil).Times(1)
		mockRepo.EXPECT().UpdateLastLogin(gomock.Any(), gomock.Any()).Return(errors.New("db error")).Times(1)

		resp, err := uc.Login(context.Background(), req)
		assert.NoError(t, err) // Login should still succeed even if last login update fails
		assert.NotNil(t, resp)
	})
}

// TestAuthUseCase_RefreshToken tests the RefreshToken method of AuthUseCase
func TestAuthUseCase_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "testsecret",
			ExpiresIn: "1h",
		},
	}
	uc := NewAuthUseCase(mockRepo, cfg)

	// Generate a valid refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
	})
	refreshTokenString, _ := refreshToken.SignedString([]byte(cfg.JWT.Secret))

	req := &models.RefreshTokenRequest{
		RefreshToken: refreshTokenString,
	}

	// Test case 1: Successful token refresh
	mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&entities.User{
		ID:              uuid.MustParse("a1b2c3d4-e5f6-7890-1234-567890abcdef"),
		Email:           "test@example.com",
		FirstName:       "John",
		LastName:        "Doe",
		IsActive:        true,
		IsEmailVerified: true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil).Times(1)
	resp, err := uc.RefreshToken(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)

	// Test case 2: Invalid refresh token
	req.RefreshToken = "invalidtoken"
	resp, err = uc.RefreshToken(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid refresh token")
}

// TestAuthUseCase_GetProfile tests the GetProfile method of AuthUseCase
func TestAuthUseCase_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "testsecret",
			ExpiresIn: "1h",
		},
	}
	uc := NewAuthUseCase(mockRepo, cfg)

	userID := uuid.New()
	expectedUser := &entities.User{
		ID:              userID,
		Email:           "test@example.com",
		FirstName:       "John",
		LastName:        "Doe",
		IsActive:        true,
		IsEmailVerified: true,
		LastLoginAt:     nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Test case 1: Successful profile retrieval
	mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(expectedUser, nil).Times(1)

	profile, err := uc.GetProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, expectedUser.Email, profile.Email)

	// Test case 2: User not found
	mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(nil, errors.New("user not found")).Times(1)

	profile, err = uc.GetProfile(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.EqualError(t, err, "user not found")
}

// TestAuthUseCase_ChangePassword tests the ChangePassword method of AuthUseCase
func TestAuthUseCase_ChangePassword(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "testsecret",
			ExpiresIn: "1h",
		},
	}

	t.Run("Successful password change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		userID := uuid.New()
		req := &models.ChangePasswordRequest{
			CurrentPassword: "oldpassword",
			NewPassword:     "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		user := &entities.User{
			ID:              userID,
			Email:           "test@example.com",
			PasswordHash:    hashPassword("oldpassword"), // Mock hashed password for "oldpassword"
			FirstName:       "John",
			LastName:        "Doe",
			IsActive:        true,
			IsEmailVerified: true,
			LastLoginAt:     nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil).Times(1)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(user, nil).Times(1)

		err := uc.ChangePassword(context.Background(), userID, req)
		assert.NoError(t, err)
	})

	t.Run("User not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		userID := uuid.New()
		req := &models.ChangePasswordRequest{
			CurrentPassword: "oldpassword",
			NewPassword:     "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(nil, errors.New("user not found")).Times(1)

		err := uc.ChangePassword(context.Background(), userID, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("Incorrect current password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		userID := uuid.New()
		req := &models.ChangePasswordRequest{
			CurrentPassword: "wrongpassword",
			NewPassword:     "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		user := &entities.User{
			ID:              userID,
			Email:           "test@example.com",
			PasswordHash:    hashPassword("oldpassword"), // Correct hashed password
			FirstName:       "John",
			LastName:        "Doe",
			IsActive:        true,
			IsEmailVerified: true,
			LastLoginAt:     nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil).Times(1)

		err := uc.ChangePassword(context.Background(), userID, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "current password is incorrect")
	})

	t.Run("Repository update error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepository(ctrl)
		uc := NewAuthUseCase(mockRepo, cfg)

		userID := uuid.New()
		req := &models.ChangePasswordRequest{
			CurrentPassword: "oldpassword",
			NewPassword:     "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		user := &entities.User{
			ID:              userID,
			Email:           "test@example.com",
			PasswordHash:    hashPassword("oldpassword"),
			FirstName:       "John",
			LastName:        "Doe",
			IsActive:        true,
			IsEmailVerified: true,
			LastLoginAt:     nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil).Times(1)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)

		err := uc.ChangePassword(context.Background(), userID, req)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to update password")
	})
}

// Helper function to hash passwords for tests
func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
