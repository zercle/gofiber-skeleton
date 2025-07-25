package usecases

import (
	"context"
	"gofiber-skeleton/internal/entities"

	"github.com/google/uuid"
)

// URLUseCase defines the interface for URL-related business logic.
type URLUseCase interface {
	CreateShortURL(ctx context.Context, originalURL string, userID uuid.UUID, customShortCode string) (*entities.URL, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]*entities.URL, error)
	UpdateShortURL(ctx context.Context, userID, urlID uuid.UUID, newOriginalURL string) (*entities.URL, error)
	DeleteShortURL(ctx context.Context, userID, urlID uuid.UUID) error
	GenerateQRCode(ctx context.Context, shortCode string) ([]byte, error)
}
