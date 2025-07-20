package usecases

import (
	"context"

	"gofiber-skeleton/internal/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, password string) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
}

type URLRepository interface {
	CreateURL(ctx context.Context, userID *string, shortCode, originalURL string) (*entities.URL, error)
	GetURLByShortCode(ctx context.Context, shortCode string) (*entities.URL, error)
	GetURLsByUserID(ctx context.Context, userID string) ([]*entities.URL, error)
	DeleteURL(ctx context.Context, id, userID string) error
}
