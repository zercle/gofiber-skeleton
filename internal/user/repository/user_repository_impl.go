package repository

import (
	"context"
	"gofiber-skeleton/internal/platform/db"
	"gofiber-skeleton/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// NewSQLUserRepository creates a new instance of SQLUserRepository.
func NewSQLUserRepository(queries *db.Queries) UserRepository {
	return &SQLUserRepository{queries: queries}
}

// SQLUserRepository implements the user.UserRepository interface and manages
// database operations for users.
type SQLUserRepository struct {
	queries *db.Queries
}

// CreateUser inserts a new user record into the database.
func (r *SQLUserRepository) CreateUser(ctx context.Context, usr *user.ModelUser) error {
	_, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username: usr.Username,
		Password: usr.Password,
		Role:     usr.Role,
	})

	return err
}

// GetUserByID retrieves a user model by their unique ID.
func (r *SQLUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*user.ModelUser, error) {
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(id)

	usr, err := r.queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		return nil, err
	}

	return &user.ModelUser{
		ID:        usr.ID.Bytes,
		Username:  usr.Username,
		Password:  usr.Password,
		Role:      usr.Role,
		CreatedAt: usr.CreatedAt.Time,
		UpdatedAt: usr.UpdatedAt.Time,
	}, nil
}

// GetUserByUsername retrieves a user model by their username.
func (r *SQLUserRepository) GetUserByUsername(ctx context.Context, username string) (*user.ModelUser, error) {
	usr, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &user.ModelUser{
		ID:        usr.ID.Bytes,
		Username:  usr.Username,
		Password:  usr.Password,
		Role:      usr.Role,
		CreatedAt: usr.CreatedAt.Time,
		UpdatedAt: usr.UpdatedAt.Time,
	}, nil
}

// UpdateUserRole updates a user's role in the database.
func (r *SQLUserRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role string) error {
	pgUUID := pgtype.UUID{}
	pgUUID.Scan(id)

	_, err := r.queries.UpdateUserRole(ctx, db.UpdateUserRoleParams{
		ID:   pgUUID,
		Role: role,
	})
	return err
}
