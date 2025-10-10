package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go -package=mocks

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) (*entity.User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates a user's information
	Update(ctx context.Context, user *entity.User) (*entity.User, error)

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a paginated list of users
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error

	// Verify marks a user as verified
	Verify(ctx context.Context, id uuid.UUID) error
}
