package repository

import (
	"context"
	"gofiber-skeleton/internal/platform/db"
	"gofiber-skeleton/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// NewSQLUserRepository creates a new instance of SQLUserRepository.
//
// Parameters:
//   - queries: Database query interface for user operations.
//
// Returns:
//   - UserRepository: An implementation of the UserRepository interface.
//
// Note:
//   This repository handles user data persistence.
func NewSQLUserRepository(queries *db.Queries) UserRepository {
	return &SQLUserRepository{queries: queries}
}

// SQLUserRepository implements the user.UserRepository interface and manages
// database operations for users.
type SQLUserRepository struct {
	queries *db.Queries
}

// CreateUser inserts a new user record into the database.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation.
//   - usr: The user model object containing user details.
//
// Returns:
//   - error: An error object if the operation fails, otherwise nil.
func (r *SQLUserRepository) CreateUser(ctx context.Context, usr *user.ModelUser) error {
	_, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username: usr.Username,
		Password: usr.Password,
		Role:     usr.Role,
	})

	return err
}

// GetUserByID retrieves a user model by their unique ID.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation.
//   - id: UUID of the user to retrieve.
//
// Returns:
//   - *user.ModelUser: User model corresponding to the ID.
//   - error: Error if retrieval fails.
func (r *SQLUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*user.ModelUser, error) {
	usr, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &user.ModelUser{ID: usr.ID.Bytes, Username: usr.Username, Password: usr.Password, Role: usr.Role, CreatedAt: usr.CreatedAt.Time, UpdatedAt: usr.UpdatedAt.Time}, nil
}

// GetUserByUsername retrieves a user model by their username.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation.
//   - username: The username string to lookup.
//
// Returns:
//   - *user.ModelUser: User model corresponding to the username.
//   - error: Error if retrieval fails.
func (r *SQLUserRepository) GetUserByUsername(ctx context.Context, username string) (*user.ModelUser, error) {
	usr, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &user.ModelUser{ID: usr.ID.Bytes, Username: usr.Username, Password: usr.Password, Role: usr.Role, CreatedAt: usr.CreatedAt.Time, UpdatedAt: usr.UpdatedAt.Time}, nil
}

// UpdateUserRole updates a user's role in the database.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation.
//   - id: UUID of the user to update.
//   - role: The new role to assign to the user.
//
// Returns:
//   - error: An error object if the operation fails, otherwise nil.
func (r *SQLUserRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error {
	_, err := r.queries.UpdateUserRole(ctx, db.UpdateUserRoleParams{
		ID:   pgtype.UUID{Bytes: id, Valid: true},
		Role: role,
	})
	return err
}
