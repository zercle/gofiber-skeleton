package usecases

import (
	"context"
	"testing"

	"gofiber-skeleton/internal/entities"
	"gofiber-skeleton/mocks"

	"go.uber.org/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestURLUseCase_CreateShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLRepo := mocks.NewMockURLRepository(ctrl)
	urlUseCase := NewURLUseCase(mockURLRepo)

	originalURL := "https://example.com"
	userID := uuid.New()

	mockURLRepo.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(nil)

	url, err := urlUseCase.CreateShortURL(context.Background(), originalURL, userID, "")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, originalURL, url.OriginalURL)
	assert.NotEmpty(t, url.ShortCode)
}

func TestURLUseCase_GetOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLRepo := mocks.NewMockURLRepository(ctrl)
	urlUseCase := NewURLUseCase(mockURLRepo)

	shortCode := "test-code"
	originalURL := "https://example.com"

	mockURLRepo.EXPECT().GetURLByShortCode(gomock.Any(), shortCode).Return(&entities.URL{OriginalURL: originalURL}, nil)

	resultURL, err := urlUseCase.GetOriginalURL(context.Background(), shortCode)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, resultURL)
}
