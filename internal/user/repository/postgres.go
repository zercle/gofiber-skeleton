package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/db"
	"github.com/zercle/gofiber-skeleton/internal/user/entity"
)

//go:generate mockgen -source=postgres.go -destination=mocks/repository.go -package=mocks

type UserRepository interface {
	Create(ctx context.Context, user *entity.User, passwordHash string) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, string, error)
}

type postgresUserRepository struct {
	queries *db.Queries
}

func NewPostgresUserRepository(queries *db.Queries) UserRepository {
	return &postgresUserRepository{queries: queries}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *entity.User, passwordHash string) error {
	params := db.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: passwordHash,
	}
	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	// Update the user entity with generated values
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time
	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, string, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	return &entity.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, dbUser.PasswordHash, nil
}
