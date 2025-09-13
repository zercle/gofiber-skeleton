package auth_repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/entities"
	auth_repositories "github.com/zercle/gofiber-skeleton/internal/domains/auth/repositories"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	sqldb "github.com/zercle/gofiber-skeleton/internal/infrastructure/database/queries"
)

type userRepository struct {
	q sqldb.Querier
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *database.Database) auth_repositories.UserRepository {
	return &userRepository{q: sqldb.New(db.Pool)}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	params := sqldb.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}
	createdUser, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.mapSQLCUserToEntity(createdUser), nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	sqlcID := pgtype.UUID{Bytes: id, Valid: true}
	user, err := r.q.GetUserByID(ctx, sqlcID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return r.mapSQLCUserToEntity(user), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return r.mapSQLCUserToEntity(user), nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	params := sqldb.UpdateUserParams{
		ID:              pgtype.UUID{Bytes: user.ID, Valid: true},
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		IsEmailVerified: pgtype.Bool{Bool: user.IsEmailVerified, Valid: true},
	}
	updatedUser, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return r.mapSQLCUserToEntity(updatedUser), nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	sqlcID := pgtype.UUID{Bytes: userID, Valid: true}
	err := r.q.UpdateUserLastLogin(ctx, sqlcID)
	if err != nil {
		return fmt.Errorf("failed to update user last login: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sqlcID := pgtype.UUID{Bytes: id, Valid: true}
	err := r.q.DeactivateUser(ctx, sqlcID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	params := sqldb.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	users, err := r.q.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var userEntities []*entities.User
	for _, user := range users {
		userEntities = append(userEntities, r.mapSQLCUserToEntity(user))
	}
	return userEntities, nil
}

func (r *userRepository) mapSQLCUserToEntity(sqlcUser sqldb.Users) *entities.User {
	return &entities.User{
		ID:              uuid.UUID(sqlcUser.ID.Bytes),
		Email:           sqlcUser.Email,
		PasswordHash:    sqlcUser.PasswordHash,
		FirstName:       sqlcUser.FirstName,
		LastName:        sqlcUser.LastName,
		IsActive:        sqlcUser.IsActive.Bool,
		IsEmailVerified: sqlcUser.IsEmailVerified.Bool,
		LastLoginAt:     r.pgtypeTimestampToTimePtr(sqlcUser.LastLoginAt),
		CreatedAt:       sqlcUser.CreatedAt.Time,
		UpdatedAt:       sqlcUser.UpdatedAt.Time,
	}
}

func (r *userRepository) pgtypeTimestampToTimePtr(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
