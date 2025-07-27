//go:generate mockgen -source=user_repository.go -destination=mocks/mock_user_repository.go -package=mocks UserRepository

//go:generate mockgen -source=user_repository.go -destination=mocks/mock_user_repository.go -package=mocks UserRepository

package repository

import (
	"context"
	"gofiber-skeleton/internal/user" // Updated import

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	CreateUser(ctx context.Context, user *user.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
}
