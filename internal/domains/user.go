package domains

import (
	"context"
	"time"
)

//go:generate mockgen -source=user.go -destination=../../test/unit/mocks/mock_user_repository.go -package=mocks

// UserRepository defines the contract for user data access
type UserRepository interface {
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Create creates a new user in the database
	Create(ctx context.Context, user *User) error

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of users
	List(ctx context.Context, limit, offset int32) ([]*User, error)

	// Count returns the total number of active users
	Count(ctx context.Context) (int64, error)
}

// UserUsecase defines the contract for user business logic
type UserUsecase interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// RegisterUser creates a new user account
	RegisterUser(ctx context.Context, input *RegisterUserInput) (*User, error)

	// UpdateUser updates user information
	UpdateUser(ctx context.Context, id string, input *UpdateUserInput) (*User, error)

	// DeleteUser deletes a user account
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves paginated user list
	ListUsers(ctx context.Context, limit, offset int32) ([]*User, error)
}

// User represents a user entity in the domain
type User struct {
	ID           string
	Email        string
	PasswordHash string
	FirstName    *string
	LastName     *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// RegisterUserInput is the input for user registration
type RegisterUserInput struct {
	Email     string
	Password  string
	FirstName *string
	LastName  *string
}

// UpdateUserInput is the input for updating a user
type UpdateUserInput struct {
	Email     *string
	FirstName *string
	LastName  *string
	IsActive  *bool
}
