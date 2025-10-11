package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
)

//go:generate mockgen -source=usecase.go -destination=../mocks/usecase_mock.go -package=mocks

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	// Register registers a new user
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.User, error)

	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetProfile retrieves the current user's profile
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)

	// Update updates a user's profile
	Update(ctx context.Context, userID uuid.UUID, req *entity.UpdateUserRequest) (*entity.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID uuid.UUID, req *entity.ChangePasswordRequest) error

	// List retrieves a paginated list of users
	List(ctx context.Context, page, perPage int) (*entity.UserListResponse, error)

	// Deactivate deactivates a user
	Deactivate(ctx context.Context, userID uuid.UUID) error

	// Activate activates a user
	Activate(ctx context.Context, userID uuid.UUID) error

	// Delete permanently deletes a user
	Delete(ctx context.Context, userID uuid.UUID) error
}
