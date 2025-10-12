package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
)

//go:generate mockgen -source=user.go -destination=../mocks/user_repository_mock.go -package=mocks

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.DomainUser) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.DomainUser, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.DomainUser, error)

	// Update updates a user
	Update(ctx context.Context, user *entity.DomainUser) error

	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error

	// Deactivate deactivates a user
	Deactivate(ctx context.Context, id uuid.UUID) error

	// List retrieves a list of users with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.DomainUser, error)

	// EmailExists checks if email already exists
	EmailExists(ctx context.Context, email string) (bool, error)

	// Count returns total number of active users
	Count(ctx context.Context) (int, error)
}
