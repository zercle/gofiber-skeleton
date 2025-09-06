package usecases

import (
	"context"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/entities"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/models"
)

type AuthUsecase interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	GetProfile(ctx context.Context, userID string) (*models.UserResponse, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.UserResponse, error)
	ChangePassword(ctx context.Context, userID string, req models.ChangePasswordRequest) error
	ListUsers(ctx context.Context, page, pageSize int) (*models.UsersListResponse, error)
	GetUser(ctx context.Context, userID string) (*models.UserResponse, error)
	DeactivateUser(ctx context.Context, userID string) error
	ActivateUser(ctx context.Context, userID string) error
	DeleteUser(ctx context.Context, userID string) error
}