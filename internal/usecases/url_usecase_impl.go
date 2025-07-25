package usecases

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"gofiber-skeleton/internal/entities"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// NewURLUseCase creates a new URLUseCase.
func NewURLUseCase(urlRepo URLRepository) URLUseCase {
	return &urlUseCase{urlRepo: urlRepo}
}

type urlUseCase struct {
	urlRepo URLRepository
}

func (uc *urlUseCase) CreateShortURL(ctx context.Context, originalURL string, userID uuid.UUID, customShortCode string) (*entities.URL, error) {
	shortCode := customShortCode
	if shortCode == "" {
		var err error
		shortCode, err = uc.generateShortCode()
		if err != nil {
			return nil, err
		}
	}

	url := &entities.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
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
	return url.OriginalURL, nil
}

func (uc *urlUseCase) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]*entities.URL, error) {
	return uc.urlRepo.GetURLsByUserID(ctx, userID)
}

func (uc *urlUseCase) UpdateShortURL(ctx context.Context, userID, urlID uuid.UUID, newOriginalURL string) (*entities.URL, error) {
	// In a real application, you would first verify that the user owns the URL.
	url := &entities.URL{
		ID:          urlID,
		OriginalURL: newOriginalURL,
	}
	err := uc.urlRepo.UpdateURL(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (uc *urlUseCase) DeleteShortURL(ctx context.Context, userID, urlID uuid.UUID) error {
	// In a real application, you would first verify that the user owns the URL.
	return uc.urlRepo.DeleteURL(ctx, urlID)
}

func (uc *urlUseCase) GenerateQRCode(ctx context.Context, shortCode string) ([]byte, error) {
	// This should be the full URL, e.g., https://your-domain.com/shortCode
	return qrcode.Encode(shortCode, qrcode.Medium, 256)
}

func (uc *urlUseCase) generateShortCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
