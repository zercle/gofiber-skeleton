package usecases

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gofiber-skeleton/internal/entities"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// NewURLUseCase creates a new URLUseCase.
func NewURLUseCase(urlRepo URLRepository, appDomain string) URLUseCase {
	return &urlUseCase{urlRepo: urlRepo, appDomain: appDomain}
}

type urlUseCase struct {
	urlRepo   URLRepository
	appDomain string
}

func (uc *urlUseCase) CreateShortURL(ctx context.Context, originalURL string, userID uuid.UUID, customShortCode string) (*entities.URL, error) {
	if customShortCode != "" {
		// Check if custom short code already exists
		_, err := uc.urlRepo.GetURLByShortCode(ctx, customShortCode)
		if err == nil {
			return nil, errors.New("custom short code already exists")
		}
	}

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
	existingURL, err := uc.urlRepo.GetURLByID(ctx, urlID)
	if err != nil {
		return nil, err
	}
	if existingURL.UserID != userID {
		return nil, errors.New("unauthorized: user does not own this URL")
	}

	url := &entities.URL{
		ID:          urlID,
		OriginalURL: newOriginalURL,
	}
	err = uc.urlRepo.UpdateURL(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (uc *urlUseCase) DeleteShortURL(ctx context.Context, userID, urlID uuid.UUID) error {
	existingURL, err := uc.urlRepo.GetURLByID(ctx, urlID)
	if err != nil {
		return err
	}
	if existingURL.UserID != userID {
		return errors.New("unauthorized: user does not own this URL")
	}
	return uc.urlRepo.DeleteURL(ctx, urlID)
}

func (uc *urlUseCase) GenerateQRCode(ctx context.Context, shortCode string) ([]byte, error) {
	fullURL := uc.appDomain + "/" + shortCode
	return qrcode.Encode(fullURL, qrcode.Medium, 256)
}

func (uc *urlUseCase) generateShortCode(ctx context.Context) (string, error) {
	for i := 0; i < 5; i++ { // Retry up to 5 times to find a unique code
		b := make([]byte, 6)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		shortCode := base64.URLEncoding.EncodeToString(b)

		// Check if the code already exists
		_, err = uc.urlRepo.GetURLByShortCode(ctx, shortCode)
		if err != nil { // Assuming "not found" is the error we want
			return shortCode, nil
		}
	}
	return "", errors.New("failed to generate a unique short code after multiple attempts")
}
