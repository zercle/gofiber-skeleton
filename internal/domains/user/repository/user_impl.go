package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/pkg/database"
)

type sqlcUserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository implementation
func NewUserRepository(injector do.Injector) (UserRepository, error) {
	db := do.MustInvoke[*database.Database](injector)
	return &sqlcUserRepository{
		db: db.GetDB(),
	}, nil
}

// Create creates a new user
func (r *sqlcUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *sqlcUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *sqlcUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user
func (r *sqlcUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email = $2, full_name = $3, updated_at = $4
		WHERE id = $1 AND is_active = true
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.FullName,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

// UpdatePassword updates user password
func (r *sqlcUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, updated_at = NOW()
		WHERE id = $1 AND is_active = true
	`

	result, err := r.db.ExecContext(ctx, query, id, passwordHash)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

// Deactivate deactivates a user
func (r *sqlcUserRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *sqlcUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT id, email, full_name, is_active, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var users []*entity.User
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// EmailExists checks if email already exists
func (r *sqlcUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE email = $1 AND is_active = true
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, email)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Count returns total number of active users
func (r *sqlcUserRepository) Count(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE is_active = true
	`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}