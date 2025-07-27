package repository

import (
	"context"
	"gofiber-skeleton/internal/platform/db"
	"gofiber-skeleton/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// NewSQLUserRepository creates a new UserRepository.
func NewSQLUserRepository(queries *db.Queries) UserRepository {
	return &SQLUserRepository{queries: queries}
}

// SQLUserRepository implements the user.UserRepository interface.
type SQLUserRepository struct {
	queries *db.Queries
}

// CreateUser creates a new user in the database.
func (r *SQLUserRepository) CreateUser(ctx context.Context, usr *user.User) error {
	_, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username: usr.Username,
		Password: usr.Password,
	})

	return err
}

// GetUserByID retrieves a user from the database by their ID.
func (r *SQLUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	usr, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &user.User{ID: usr.ID.Bytes, Username: usr.Username, Password: usr.Password, CreatedAt: usr.CreatedAt.Time, UpdatedAt: usr.UpdatedAt.Time}, nil
}

// GetUserByUsername retrieves a user from the database by their username.
func (r *SQLUserRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	usr, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &user.User{ID: usr.ID.Bytes, Username: usr.Username, Password: usr.Password, CreatedAt: usr.CreatedAt.Time, UpdatedAt: usr.UpdatedAt.Time}, nil
}