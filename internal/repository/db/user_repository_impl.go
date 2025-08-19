package db

import (
	"context"
	"database/sql"
	"github.com/zercle/gofiber-skeleton/internal/domain"
)

type userRepo struct {
	queries *Queries
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepo{
		queries: New(db),
	}
}

func (ur *userRepo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	createdUser, err := ur.queries.CreateUser(ctx, CreateUserParams{
		Username:    user.Username,
		PasswordHash: user.PasswordHash,
		Role:        user.Role,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           int64(createdUser.ID),
		Username:     createdUser.Username,
		PasswordHash: createdUser.PasswordHash,
		Role:         createdUser.Role,
		CreatedAt:    createdUser.CreatedAt,
	}, nil
}

func (ur *userRepo) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := ur.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           int64(user.ID),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func (ur *userRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := ur.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           int64(user.ID),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}, nil
}