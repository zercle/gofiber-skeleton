package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
)

//go:generate mockgen -source=repository.go -destination=../mocks/repository_mock.go -package=mocks

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// List retrieves a paginated list of users
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)

	// Count returns the total number of active users
	Count(ctx context.Context) (int64, error)

	// Update updates a user
	Update(ctx context.Context, user *entity.User) error

	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error

	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	// Verify marks a user as verified
	Verify(ctx context.Context, userID uuid.UUID) error

	// Deactivate deactivates a user
	Deactivate(ctx context.Context, userID uuid.UUID) error

	// Activate activates a user
	Activate(ctx context.Context, userID uuid.UUID) error

	// Delete permanently deletes a user
	Delete(ctx context.Context, userID uuid.UUID) error
}
