//go:generate mockgen -source=url_repository.go -destination=mocks/mock_url_repository.go -package=mocks URLRepository

package usecases

import (
	"context"
	"gofiber-skeleton/internal/entities"

	"github.com/google/uuid"
)

// URLRepository defines the interface for URL data operations.
type URLRepository interface {
	CreateURL(ctx context.Context, url *entities.URL) (*entities.URL, error)
	GetURLByShortCode(ctx context.Context, shortCode string) (*entities.URL, error)
	GetURLByID(ctx context.Context, id uuid.UUID) (*entities.URL, error) // Added for ownership verification
	GetURLsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.URL, error)
	UpdateURL(ctx context.Context, url *entities.URL) (*entities.URL, error) // Updated to return *entities.URL
	DeleteURL(ctx context.Context, id uuid.UUID) error
}
