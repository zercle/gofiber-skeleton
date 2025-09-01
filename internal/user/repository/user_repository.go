package repository

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
	"github.com/zercle/gofiber-skeleton/internal/user"
)

// userRepository implements the user.UserRepository interface.
type userRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(queries *sqlc.Queries) user.UserRepository {
	return &userRepository{
		queries: queries,
	}
}

// CreateUser creates a new user in the database.
func (r *userRepository) CreateUser(u user.User) (user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := sqlc.CreateUserParams{
		Email:    u.Email,
		Password: u.Password,
	}

	createdUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return user.User{}, err
	}

	return toDomainUser(createdUser), nil
}

func (r *userRepository) GetUserByEmail(ctx *fiber.Ctx, email string) (user.User, error) {
	foundUser, err := r.queries.GetUserByEmail(ctx.Context(), email)
	if err != nil {
		return user.User{}, err
	}

	return toDomainUser(foundUser), nil
}

// toDomainUser maps a sqlc.User to a user.User domain model.
func toDomainUser(u sqlc.User) user.User {
	return user.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}