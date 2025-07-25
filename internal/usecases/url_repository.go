package usecases

import (
	"context"
	"gofiber-skeleton/internal/entities"

	"github.com/google/uuid"
)

// URLRepository defines the interface for URL data operations.
type URLRepository interface {
	CreateURL(ctx context.Context, url *entities.URL) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*entities.URL, error)
	GetURLsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.URL, error)
	UpdateURL(ctx context.Context, url *entities.URL) error
	DeleteURL(ctx context.Context, id uuid.UUID) error
}
