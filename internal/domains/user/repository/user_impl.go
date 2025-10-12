package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/pkg/database"
)

type sqlcUserRepository struct {
	db      *database.Database
	queries *entity.Queries
}

// NewUserRepository creates a new user repository implementation
func NewUserRepository(injector do.Injector) (UserRepository, error) {
	db := do.MustInvoke[*database.Database](injector)
	return &sqlcUserRepository{
		db:      db,
		queries: entity.New(db.GetPool()),
	}, nil
}

// Create creates a new user
func (r *sqlcUserRepository) Create(ctx context.Context, user *entity.DomainUser) error {
	params := entity.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
	}

	_, err := r.queries.CreateUser(ctx, params)
	return err
}

// GetByID retrieves a user by ID
func (r *sqlcUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.DomainUser, error) {
	var userUUID [16]byte
	copy(userUUID[:], id[:])
	pgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	user, err := r.queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	domainUser := &entity.DomainUser{
		ID:           uuid.UUID(user.ID.Bytes),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		IsActive:     user.IsActive.Bool,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}

	return domainUser, nil
}

// GetByEmail retrieves a user by email
func (r *sqlcUserRepository) GetByEmail(ctx context.Context, email string) (*entity.DomainUser, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	domainUser := &entity.DomainUser{
		ID:           uuid.UUID(user.ID.Bytes),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		IsActive:     user.IsActive.Bool,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}

	return domainUser, nil
}

// Update updates a user
func (r *sqlcUserRepository) Update(ctx context.Context, user *entity.DomainUser) error {
	var userUUID [16]byte
	copy(userUUID[:], user.ID[:])
	pgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	params := entity.UpdateUserParams{
		ID:       pgUUID,
		Email:    user.Email,
		FullName: user.FullName,
	}

	_, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.ErrUserNotFound
		}
		return err
	}

	return nil
}

// UpdatePassword updates user password
func (r *sqlcUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	var userUUID [16]byte
	copy(userUUID[:], id[:])
	pgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	params := entity.UpdateUserPasswordParams{
		ID:           pgUUID,
		PasswordHash: passwordHash,
	}

	_, err := r.queries.UpdateUserPassword(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.ErrUserNotFound
		}
		return err
	}

	return nil
}

// Deactivate deactivates a user
func (r *sqlcUserRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	var userUUID [16]byte
	copy(userUUID[:], id[:])
	pgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	_, err := r.queries.DeactivateUser(ctx, pgUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.ErrUserNotFound
		}
		return err
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *sqlcUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.DomainUser, error) {
	params := entity.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	var users []*entity.DomainUser
	for _, row := range rows {
		user := &entity.DomainUser{
			ID:           uuid.UUID(row.ID.Bytes),
			Email:        row.Email,
			PasswordHash: "", // Not included in list query
			FullName:     row.FullName,
			IsActive:     row.IsActive.Bool,
			CreatedAt:    row.CreatedAt.Time,
			UpdatedAt:    row.UpdatedAt.Time,
		}
		users = append(users, user)
	}

	return users, nil
}

// EmailExists checks if email already exists
func (r *sqlcUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	count, err := r.queries.UserExists(ctx, email)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Count returns total number of active users
func (r *sqlcUserRepository) Count(ctx context.Context) (int, error) {
	count, err := r.queries.UserExists(ctx, "")
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
