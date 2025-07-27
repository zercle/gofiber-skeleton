package repository

import (
	"context"
	"gofiber-skeleton/internal/entities"
	db "gofiber-skeleton/internal/repository/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// NewUserRepository creates a new UserRepository.
func NewUserRepository(querier DBQueriesInterface) *UserRepository {
	return &UserRepository{queries: querier}
}

// UserRepository implements the usecases.UserRepository interface.
type UserRepository struct {
	queries DBQueriesInterface
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	_, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
	})
	return err
}

// GetUserByID retrieves a user from the database by their ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &entities.User{ID: user.ID.Bytes, Username: user.Username, Password: user.Password, CreatedAt: user.CreatedAt.Time, UpdatedAt: user.UpdatedAt.Time}, nil
}

// GetUserByUsername retrieves a user from the database by their username.
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &entities.User{ID: user.ID.Bytes, Username: user.Username, Password: user.Password, CreatedAt: user.CreatedAt.Time, UpdatedAt: user.UpdatedAt.Time}, nil
}
