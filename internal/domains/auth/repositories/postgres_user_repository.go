package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/entities"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	_, err := r.db.Exec(ctx, query, 
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.IsActive, user.CreatedAt, user.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`
	
	user := &entities.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL`
	
	user := &entities.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $2, password = $3, first_name = $4, last_name = $5, 
		    is_active = $6, updated_at = $7
		WHERE id = $1 AND deleted_at IS NULL`
	
	result, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.IsActive, user.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return types.ErrUserNotFound
	}
	
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users 
		SET deleted_at = NOW(), updated_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return types.ErrUserNotFound
	}
	
	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()
	
	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
			&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (r *PostgresUserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	return count, nil
}

func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`
	
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists by email: %w", err)
	}
	
	return exists, nil
}