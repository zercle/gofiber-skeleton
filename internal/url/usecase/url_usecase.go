//go:generate mockgen -source=url_usecase.go -destination=mocks/mock_url_usecase.go -package=mocks URLUseCase

package usecase

import (
	"context"
	"gofiber-skeleton/internal/url"
)

// URLUseCase defines the interface for URL-related business logic.
type URLUseCase interface {
	CreateURL(ctx context.Context, originalURL string, userID string) (*url.ModelURL, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}
