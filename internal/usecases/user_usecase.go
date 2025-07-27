//go:generate mockgen -source=user_usecase.go -destination=mocks/mock_user_usecase.go -package=mocks UserUseCase

package usecases

import (
	"context"
	"gofiber-skeleton/internal/entities"
)

// UserUseCase defines the interface for user-related business logic.
type UserUseCase interface {
	Register(ctx context.Context, username, password string) (*entities.User, error)
	Login(ctx context.Context, username, password string) (string, error) // Returns JWT
}
