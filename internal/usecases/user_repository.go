//go:generate mockgen -source=user_repository.go -destination=mocks/mock_user_repository.go -package=mocks UserRepository

package usecases

import (
	"context"
	"gofiber-skeleton/internal/entities"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
}
