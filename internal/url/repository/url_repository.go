//go:generate mockgen -source=url_repository.go -destination=mocks/mock_url_repository.go -package=mocks URLRepository

package repository

import (
	"context"
	"gofiber-skeleton/internal/url" // Updated import

	"github.com/google/uuid"
)

// URLRepository defines the interface for URL data access.
type URLRepository interface {
	CreateURL(ctx context.Context, url *url.URL) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*url.URL, error)
	GetURLByID(ctx context.Context, id uuid.UUID) (*url.URL, error)
}
