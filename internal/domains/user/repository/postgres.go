package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/db"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	queries *db.Queries
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(dbConn *sql.DB) UserRepository {
	return &PostgresUserRepository{
		queries: db.New(dbConn),
	}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.dbUserToEntity(&dbUser), nil
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return r.dbUserToEntity(&dbUser), nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return r.dbUserToEntity(&dbUser), nil
}

// Update updates a user's information
func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:         user.ID,
		FullName:   user.FullName,
		IsActive:   sql.NullBool{Bool: user.IsActive, Valid: true},
		IsVerified: sql.NullBool{Bool: user.IsVerified, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.dbUserToEntity(&dbUser), nil
}

// Delete deletes a user by ID
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves a paginated list of users
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	dbUsers, err := r.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.dbUserToEntity(&dbUser)
	}

	return users, nil
}

// Count returns the total number of users
func (r *PostgresUserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// UpdatePassword updates a user's password
func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if err := r.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	}); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// Verify marks a user as verified
func (r *PostgresUserRepository) Verify(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.VerifyUser(ctx, id); err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}
	return nil
}

// dbUserToEntity converts database user model to entity
func (r *PostgresUserRepository) dbUserToEntity(dbUser *db.User) *entity.User {
	return &entity.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		FullName:     dbUser.FullName,
		IsActive:     dbUser.IsActive.Bool,
		IsVerified:   dbUser.IsVerified.Bool,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}
}
