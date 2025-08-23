package repository

import (
	"context"
	"database/sql"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/repository/db"
)

type userRepository struct {
	db *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *db.Queries) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	ctx := context.Background()
	
	dbUser, err := r.db.CreateUser(ctx, db.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return err
	}

	user.ID = dbUser.ID.String()
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	ctx := context.Background()
	
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	dbUser, err := r.db.GetUserByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           dbUser.ID.String(),
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	ctx := context.Background()
	
	dbUser, err := r.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           dbUser.ID.String(),
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}, nil
}

func (r *userRepository) Update(user *domain.User) error {
	ctx := context.Background()
	
	uuid, err := parseUUID(user.ID)
	if err != nil {
		return err
	}

	dbUser, err := r.db.UpdateUser(ctx, db.UpdateUserParams{
		ID:           uuid,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return err
	}

	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

func (r *userRepository) Delete(id string) error {
	ctx := context.Background()
	
	uuid, err := parseUUID(id)
	if err != nil {
		return err
	}

	return r.db.DeleteUser(ctx, uuid)
}