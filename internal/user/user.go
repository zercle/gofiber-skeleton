//go:generate go run go.uber.org/mock/mockgen -source=user.go -destination=mock/user_mock.go -package=mock

// Package user provides domain models and interfaces for user management.
package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity with optimized field alignment.
// Fields are ordered by size (largest to smallest) for memory efficiency on 64-bit architectures.
type User struct {
	// 24 bytes
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z"`
	// 16 bytes
	ID       uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string    `json:"email" example:"user@example.com"`
	Password string    `json:"-"`
}

// UserRepository defines data access operations for user entities.
type UserRepository interface {
	// CreateUser persists a new user to the database.
	// It returns the created user with its ID and timestamps, or an error
	// if a user with the same email already exists.
	CreateUser(user User) (User, error)
}

// UserUsecase defines business logic operations for user management.
type UserUsecase interface {
	// Register handles the business logic for creating a new user account.
	// It validates the input payload, hashes the password, and then
	// persists the new user via the UserRepository. It returns the newly
	// created user or an error if validation fails or the user already exists.
	Register(payload RegisterPayload) (User, error)
}

// RegisterPayload contains the required data for user registration.
type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"securepassword123"`
}