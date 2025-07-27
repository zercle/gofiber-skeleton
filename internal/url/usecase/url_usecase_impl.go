package usecase

import (
	"context"
	"errors" // Added errors import
	"gofiber-skeleton/internal/url"
	"gofiber-skeleton/internal/url/repository"
	"time"

	"gofiber-skeleton/pkg/utils" // Assuming a utility for short code generation

	"github.com/google/uuid" // Added uuid import
)

// NewURLUseCase creates a new URLUseCase.
func NewURLUseCase(urlRepo repository.URLRepository, cacheRepo any) URLUseCase { // cacheRepo is a placeholder for now
	return &urlUseCase{urlRepo: urlRepo}
}

type urlUseCase struct {
	urlRepo repository.URLRepository
	// cacheRepo cache.CacheRepository // Uncomment when cache is integrated
}

func (uc *urlUseCase) CreateURL(ctx context.Context, originalURL string, userID string) (*url.URL, error) {
	// Generate short code
	shortCode := utils.GenerateShortCode(8) // Assuming a utility function

	// Convert userID string to uuid.UUID
	var userUUID uuid.UUID
	if userID != "" {
		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			return nil, err
		}
		userUUID = parsedUserID
	}

	// Set expiration time (e.g., 24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	url := &url.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userUUID,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}

	err := uc.urlRepo.CreateURL(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (uc *urlUseCase) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := uc.urlRepo.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	// Check if URL has expired
	if url.ExpiresAt.Before(time.Now()) {
		return "", errors.New("URL expired") // Assuming errors package
	}

	return url.OriginalURL, nil
}
