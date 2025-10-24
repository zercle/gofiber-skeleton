package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/zercle/template-go-fiber/internal/domains"
	db "github.com/zercle/template-go-fiber/internal/infrastructure/sqlc"
)

// UserRepository implements the domains.UserRepository interface using sqlc
type UserRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domains.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return r.dbUserToDomain(dbUser), nil
}

// GetByEmail retrieves a user by their email address
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domains.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return r.dbUserToDomain(dbUser), nil
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *domains.User) error {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    r.stringToNullString(user.FirstName),
		LastName:     r.stringToNullString(user.LastName),
		IsActive:     sql.NullBool{Bool: user.IsActive, Valid: true},
	})
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domains.User) error {
	return r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    r.stringToNullString(user.FirstName),
		LastName:     r.stringToNullString(user.LastName),
		IsActive:     sql.NullBool{Bool: user.IsActive, Valid: true},
	})
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteUser(ctx, id)
}

// List retrieves a paginated list of users
func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]*domains.User, error) {
	dbUsers, err := r.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	users := make([]*domains.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.dbUserToDomain(dbUser)
	}

	return users, nil
}

// Count returns the total number of active users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

// Helper functions

func (r *UserRepository) dbUserToDomain(dbUser db.User) *domains.User {
	return &domains.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		FirstName:    r.nullStringToPointer(dbUser.FirstName),
		LastName:     r.nullStringToPointer(dbUser.LastName),
		IsActive:     dbUser.IsActive.Bool,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
		DeletedAt:    r.nullTimeToPointer(dbUser.DeletedAt),
	}
}

func (r *UserRepository) stringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func (r *UserRepository) nullStringToPointer(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func (r *UserRepository) nullTimeToPointer(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
