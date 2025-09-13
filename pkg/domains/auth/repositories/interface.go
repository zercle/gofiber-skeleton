package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/entities"
)

// UserRepository defines the interface for user-related operations.
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
}
