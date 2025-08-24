package userrepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

type userRepository struct {
	db sqlc.Querier
}

// NewUserRepository creates a new user repository
func NewUserRepository(db sqlc.Querier) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	ctx := context.Background()

	dbUser, err := r.db.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return nil, err
	}

	user.ID = dbUser.ID.String()
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time
	return user, nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	dbUser, err := r.db.GetUserByID(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           dbUser.ID.String(),
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
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
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) Update(user *domain.User) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return err
	}

	dbUser, err := r.db.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:           parsedUUID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return err
	}

	user.UpdatedAt = dbUser.UpdatedAt.Time
	return nil
}

func (r *userRepository) Delete(id string) error {
	ctx := context.Background()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.DeleteUser(ctx, parsedUUID)
}
