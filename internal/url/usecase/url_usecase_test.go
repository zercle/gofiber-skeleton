package usecase_test

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/url"
	mock_repository "gofiber-skeleton/internal/url/repository/mocks"
	"gofiber-skeleton/internal/url/usecase"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestURLUseCase_CreateURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLRepo := mock_repository.NewMockURLRepository(ctrl)
	urlUseCase := usecase.NewURLUseCase(mockURLRepo, nil) // nil for cacheRepo as it's not used in this test

	ctx := context.Background()
	originalURL := "https://example.com"
	userID := uuid.New().String()

	// Test case 1: Successful URL creation
	mockURLRepo.EXPECT().CreateURL(ctx, gomock.Any()).Return(nil).Times(1)

	createdURL, err := urlUseCase.CreateURL(ctx, originalURL, userID)
	assert.NoError(t, err)
	assert.NotNil(t, createdURL)
	assert.Equal(t, originalURL, createdURL.OriginalURL)
	assert.NotEmpty(t, createdURL.ShortCode)
	assert.Equal(t, userID, createdURL.UserID.String())
	assert.WithinDuration(t, time.Now(), createdURL.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), createdURL.ExpiresAt, 5*time.Second)

	// Test case 2: Error creating URL in repository
	repoError := errors.New("failed to create URL in repository")
	mockURLRepo.EXPECT().CreateURL(ctx, gomock.Any()).Return(repoError).Times(1)

	createdURL, err = urlUseCase.CreateURL(ctx, originalURL, userID)
	assert.Error(t, err)
	assert.Nil(t, createdURL)
	assert.Equal(t, repoError, err)
}

func TestURLUseCase_GetOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLRepo := mock_repository.NewMockURLRepository(ctrl)
	urlUseCase := usecase.NewURLUseCase(mockURLRepo, nil)

	ctx := context.Background()
	shortCode := "abc123def"
	expectedURL := &url.URL{
		OriginalURL: "https://example.com/original",
		ShortCode:   shortCode,
		ExpiresAt:   time.Now().Add(time.Hour), // Not expired
	}

	// Test case 1: Successful retrieval of original URL
	mockURLRepo.EXPECT().GetURLByShortCode(ctx, shortCode).Return(expectedURL, nil).Times(1)

	retrievedURL, err := urlUseCase.GetOriginalURL(ctx, shortCode)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL.OriginalURL, retrievedURL)

	// Test case 2: Repository error during retrieval
	repoError := errors.New("URL not found")
	mockURLRepo.EXPECT().GetURLByShortCode(ctx, "nonexistent").Return(nil, repoError).Times(1)

	retrievedURL, err = urlUseCase.GetOriginalURL(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Empty(t, retrievedURL)
	assert.Equal(t, repoError, err)

	// Test case 3: Expired URL
	expiredURL := &url.URL{
		OriginalURL: "https://example.com/expired",
		ShortCode:   "expired",
		ExpiresAt:   time.Now().Add(-time.Hour), // Expired
	}
	mockURLRepo.EXPECT().GetURLByShortCode(ctx, "expired").Return(expiredURL, nil).Times(1)

	retrievedURL, err = urlUseCase.GetOriginalURL(ctx, "expired")
	assert.Error(t, err)
	assert.Empty(t, retrievedURL)
	assert.Equal(t, "URL expired", err.Error())
}
