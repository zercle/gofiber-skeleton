package usecases

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"gofiber-skeleton/internal/entities"
)

type URLUsecase struct {
	repo URLRepository
}

func NewURLUsecase(repo URLRepository) *URLUsecase {
	return &URLUsecase{repo: repo}
}

func (uc *URLUsecase) CreateShortURL(ctx context.Context, userID *string, originalURL string, customShortCode ...string) (*entities.URL, error) {
	shortCode := ""
	if len(customShortCode) > 0 && customShortCode[0] != "" {
		shortCode = customShortCode[0]
	} else {
		var err error
		shortCode, err = uc.generateShortCode()
		if err != nil {
			return nil, err
		}
	}

	return uc.repo.CreateURL(ctx, userID, shortCode, originalURL)
}

func (uc *URLUsecase) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	url, err := uc.repo.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}

func (uc *URLUsecase) GetUserURLs(ctx context.Context, userID string) ([]*entities.URL, error) {
	return uc.repo.GetURLsByUserID(ctx, userID)
}

func (uc *URLUsecase) DeleteURL(ctx context.Context, id, userID string) error {
	return uc.repo.DeleteURL(ctx, id, userID)
}

func (uc *URLUsecase) generateShortCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
