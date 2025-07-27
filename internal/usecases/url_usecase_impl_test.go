package usecases_test

import (
	"context"
	"database/sql"
	"testing"
	// "time"

	"gofiber-skeleton/internal/entities"
	usecases "gofiber-skeleton/internal/usecases"
	"gofiber-skeleton/internal/usecases/mocks"

	"go.uber.org/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestURLUseCase_CreateShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLRepo := mocks.NewMockURLRepository(ctrl)
	urlUseCase := usecases.NewURLUseCase(mockURLRepo, "http://localhost:8080")

	originalURL := "https://example.com"
	userID := uuid.New()


	mockURLRepo.EXPECT().GetURLByShortCode(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows).MinTimes(1).MaxTimes(5)
	mockURLRepo.EXPECT().CreateURL(gomock.Any(), gomock.AssignableToTypeOf(&entities.URL{})).Return(&entities.URL{
		ID:          uuid.New(),
		OriginalURL: originalURL,
		ShortCode:   "mock-short-code",
		UserID:      userID,
	}, nil)

	println("Calling CreateShortURL with originalURL:", originalURL, "userID:", userID.String())
	url, err := urlUseCase.CreateShortURL(context.Background(), originalURL, userID, "")
	println("CreateShortURL returned url:", url, "err:", err)

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, originalURL, url.OriginalURL)
	assert.NotEmpty(t, url.ShortCode)
}

func TestURLUseCase_GetOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLRepo := mocks.NewMockURLRepository(ctrl)
	urlUseCase := usecases.NewURLUseCase(mockURLRepo, "http://localhost:8080")

	shortCode := "test-code"
	originalURL := "https://example.com"

	mockURLRepo.EXPECT().GetURLByShortCode(gomock.Any(), shortCode).Return(&entities.URL{OriginalURL: originalURL}, nil)

	resultURL, err := urlUseCase.GetOriginalURL(context.Background(), shortCode)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, resultURL)
}
