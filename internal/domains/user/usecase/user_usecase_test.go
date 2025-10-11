package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/mocks"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/pkg/auth"
	"github.com/zercle/gofiber-skeleton/pkg/validator"
	"go.uber.org/mock/gomock"
)

func TestUserUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: 24 * time.Hour,
		},
	}
	jwtManager := auth.NewJWTManager(cfg)
	validatorInstance := validator.New()
	uc := usecase.NewUserUsecase(mockRepo, jwtManager, validatorInstance, cfg)

	ctx := context.Background()
	fullName := "John Doe"
	req := &entity.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "Password123",
		FullName: &fullName,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(nil, sql.ErrNoRows)

		mockRepo.EXPECT().
			GetByUsername(ctx, req.Username).
			Return(nil, sql.ErrNoRows)

		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, user *entity.User) error {
				user.ID = uuid.New()
				user.IsActive = true
				user.IsVerified = false
				user.CreatedAt = time.Now()
				user.UpdatedAt = time.Now()
				return nil
			})

		user, err := uc.Register(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if user.Username != req.Username {
			t.Errorf("Expected username %s, got %s", req.Username, user.Username)
		}
		if user.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, user.Email)
		}
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		existingUser := &entity.User{
			ID:    uuid.New(),
			Email: req.Email,
		}

		mockRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(existingUser, nil)

		_, err := uc.Register(ctx, req)
		if err == nil {
			t.Fatal("Expected error for existing email")
		}
	})

	t.Run("Username Already Exists", func(t *testing.T) {
		existingUser := &entity.User{
			ID:       uuid.New(),
			Username: req.Username,
		}

		mockRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(nil, sql.ErrNoRows)

		mockRepo.EXPECT().
			GetByUsername(ctx, req.Username).
			Return(existingUser, nil)

		_, err := uc.Register(ctx, req)
		if err == nil {
			t.Fatal("Expected error for existing username")
		}
	})
}

func TestUserUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test-secret",
			ExpiresIn:        24 * time.Hour,
			RefreshExpiresIn: 168 * time.Hour,
		},
	}
	jwtManager := auth.NewJWTManager(cfg)
	validatorInstance := validator.New()
	uc := usecase.NewUserUsecase(mockRepo, jwtManager, validatorInstance, cfg)

	ctx := context.Background()
	password := "Password123"
	passwordHash, _ := auth.HashPassword(password)

	user := &entity.User{
		ID:           uuid.New(),
		Username:     "johndoe",
		Email:        "john@example.com",
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	req := &entity.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(user, nil)

		mockRepo.EXPECT().
			UpdateLastLogin(ctx, user.ID).
			Return(nil)

		resp, err := uc.Login(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.AccessToken == "" {
			t.Error("Expected access token")
		}
		if resp.RefreshToken == "" {
			t.Error("Expected refresh token")
		}
		if resp.User.Email != user.Email {
			t.Errorf("Expected user email %s, got %s", user.Email, resp.User.Email)
		}
	})

	t.Run("Invalid Email", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(nil, sql.ErrNoRows)

		_, err := uc.Login(ctx, req)
		if err == nil {
			t.Fatal("Expected error for invalid email")
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(user, nil)

		invalidReq := &entity.LoginRequest{
			Email:    req.Email,
			Password: "WrongPassword123",
		}

		_, err := uc.Login(ctx, invalidReq)
		if err == nil {
			t.Fatal("Expected error for invalid password")
		}
	})

	t.Run("Inactive User", func(t *testing.T) {
		inactiveUser := &entity.User{
			ID:           user.ID,
			Email:        user.Email,
			PasswordHash: passwordHash,
			IsActive:     false,
		}

		mockRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(inactiveUser, nil)

		_, err := uc.Login(ctx, req)
		if err == nil {
			t.Fatal("Expected error for inactive user")
		}
	})
}

func TestUserUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: 24 * time.Hour,
		},
	}
	jwtManager := auth.NewJWTManager(cfg)
	validatorInstance := validator.New()
	uc := usecase.NewUserUsecase(mockRepo, jwtManager, validatorInstance, cfg)

	ctx := context.Background()
	userID := uuid.New()
	user := &entity.User{
		ID:       userID,
		Username: "johndoe",
		Email:    "john@example.com",
		IsActive: true,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		result, err := uc.GetByID(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.ID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, result.ID)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, sql.ErrNoRows)

		_, err := uc.GetByID(ctx, userID)
		if err == nil {
			t.Fatal("Expected error for user not found")
		}
	})
}

func TestUserUsecase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: 24 * time.Hour,
		},
	}
	jwtManager := auth.NewJWTManager(cfg)
	validatorInstance := validator.New()
	uc := usecase.NewUserUsecase(mockRepo, jwtManager, validatorInstance, cfg)

	ctx := context.Background()
	userID := uuid.New()
	user := &entity.User{
		ID:       userID,
		Username: "johndoe",
		Email:    "john@example.com",
		IsActive: true,
	}

	newFullName := "Jane Doe"
	req := &entity.UpdateUserRequest{
		FullName: &newFullName,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockRepo.EXPECT().
			Update(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, u *entity.User) error {
				u.UpdatedAt = time.Now()
				return nil
			})

		result, err := uc.Update(ctx, userID, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if *result.FullName != newFullName {
			t.Errorf("Expected full name %s, got %s", newFullName, *result.FullName)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, sql.ErrNoRows)

		_, err := uc.Update(ctx, userID, req)
		if err == nil {
			t.Fatal("Expected error for user not found")
		}
	})
}

func TestUserUsecase_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: 24 * time.Hour,
		},
	}
	jwtManager := auth.NewJWTManager(cfg)
	validatorInstance := validator.New()
	uc := usecase.NewUserUsecase(mockRepo, jwtManager, validatorInstance, cfg)

	ctx := context.Background()
	userID := uuid.New()
	currentPassword := "Password123"
	newPassword := "NewPassword123"
	passwordHash, _ := auth.HashPassword(currentPassword)

	user := &entity.User{
		ID:           userID,
		PasswordHash: passwordHash,
	}

	req := &entity.ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockRepo.EXPECT().
			UpdatePassword(ctx, userID, gomock.Any()).
			Return(nil)

		err := uc.ChangePassword(ctx, userID, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Wrong Current Password", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		wrongReq := &entity.ChangePasswordRequest{
			CurrentPassword: "WrongPassword123",
			NewPassword:     newPassword,
		}

		err := uc.ChangePassword(ctx, userID, wrongReq)
		if err == nil {
			t.Fatal("Expected error for wrong current password")
		}
	})
}

func TestUserUsecase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: 24 * time.Hour,
		},
	}
	jwtManager := auth.NewJWTManager(cfg)
	validatorInstance := validator.New()
	uc := usecase.NewUserUsecase(mockRepo, jwtManager, validatorInstance, cfg)

	ctx := context.Background()
	page := 1
	perPage := 20

	users := []*entity.User{
		{ID: uuid.New(), Username: "user1", Email: "user1@example.com"},
		{ID: uuid.New(), Username: "user2", Email: "user2@example.com"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			List(ctx, perPage, 0).
			Return(users, nil)

		mockRepo.EXPECT().
			Count(ctx).
			Return(int64(2), nil)

		result, err := uc.List(ctx, page, perPage)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(result.Users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(result.Users))
		}
		if result.TotalCount != 2 {
			t.Errorf("Expected total count 2, got %d", result.TotalCount)
		}
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.EXPECT().
			List(ctx, perPage, 0).
			Return(nil, errors.New("database error"))

		_, err := uc.List(ctx, page, perPage)
		if err == nil {
			t.Fatal("Expected error from repository")
		}
	})
}
