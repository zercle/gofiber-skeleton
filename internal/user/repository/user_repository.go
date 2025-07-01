//go:generate mockgen -source=user_repository.go -destination=../mocks/mock_user_repository.go -package=mocks
package repository

import (
	"context"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/pkg/types"
)

type UserRepository interface {
	GetUser(ctx context.Context, id types.UUIDv7) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id types.UUIDv7) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
