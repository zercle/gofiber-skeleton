package usecase

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/url"
	"gofiber-skeleton/internal/url/repository"
	"time"

	"gofiber-skeleton/pkg/utils"

	"github.com/google/uuid"
)

// NewURLUseCase creates a new instance of URLUseCase.
//
// Parameters:
//   - urlRepo: Repository interface for URL persistence.
//   - cacheRepo: Placeholder for cache repository, currently unused.
//
// Returns:
//   - URLUseCase: Implementation of URLUseCase interface.
func NewURLUseCase(urlRepo repository.URLRepository, cacheRepo any) URLUseCase {
	return &urlUseCase{urlRepo: urlRepo}
}

type urlUseCase struct {
	urlRepo repository.URLRepository
	// cacheRepo cache.CacheRepository // Uncomment when cache is integrated
}

// CreateURL generates a short code, constructs a ModelURL, and persists it.
//
// Parameters:
//   - ctx: Context for request management.
//   - originalURL: The original URL string to shorten.
//   - userID: Optional user ID string associated with the URL.
//
// Returns:
//   - *url.ModelURL: The created URL model with generated short code.
//   - error: Error if creation fails.
func (uc *urlUseCase) CreateURL(ctx context.Context, originalURL string, userID string) (*url.ModelURL, error) {
	shortCode := utils.GenerateShortCode(8)

	var userUUID uuid.UUID
	if userID != "" {
		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			return nil, err
		}
		userUUID = parsedUserID
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	url := &url.ModelURL{
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

// GetOriginalURL retrieves the original URL for a given short code.
//
// Parameters:
//   - ctx: Context for request management.
//   - shortCode: The short code string to look up.
//
// Returns:
//   - string: The original URL corresponding to the short code.
//   - error: Error if the URL is not found or expired.
func (uc *urlUseCase) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := uc.urlRepo.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if url.ExpiresAt.Before(time.Now()) {
		return "", errors.New("URL expired")
	}

	return url.OriginalURL, nil
}
