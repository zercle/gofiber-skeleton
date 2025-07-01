//go:generate mockgen -source=user_repository.go -destination=../mocks/mock_user_repository.go -package=mocks
package repository

import (
	"context"
	"gofiber-skeleton/internal/user/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
