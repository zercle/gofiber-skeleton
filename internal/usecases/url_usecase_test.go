package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"gofiber-skeleton/internal/entities"
	"gofiber-skeleton/mocks"
)

func TestURLUsecase_CreateShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	usecase := NewURLUsecase(mockRepo)

	longURL := "https://example.com"
	var userID string
	u, err := uuid.NewV7()
	assert.NoError(t, err)
	userID = u.String()

	// Test with generated short code
	u7, err := uuid.NewV7()
	assert.NoError(t, err)
	expectedURL := &entities.URL{
		ID:          u7.String(),
		UserID:      &userID,
		ShortCode:   "generated",
		OriginalURL: longURL,
		CreatedAt:   time.Now(),
	}
	mockRepo.EXPECT().CreateURL(gomock.Any(), gomock.Any(), gomock.Any(), longURL).Return(expectedURL, nil)
	url, err := usecase.CreateShortURL(context.Background(), &userID, longURL)
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, expectedURL.OriginalURL, url.OriginalURL)

	// Test with custom short code
	customShortCode := "mycustomcode"
	u7_2, err := uuid.NewV7()
	assert.NoError(t, err)
	expectedURL2 := &entities.URL{
		ID:          u7_2.String(),
		UserID:      &userID,
		ShortCode:   customShortCode,
		OriginalURL: longURL,
		CreatedAt:   time.Now(),
	}
	mockRepo.EXPECT().CreateURL(gomock.Any(), gomock.Any(), customShortCode, longURL).Return(expectedURL2, nil)
	url2, err := usecase.CreateShortURL(context.Background(), &userID, longURL, customShortCode)
	assert.NoError(t, err)
	assert.NotNil(t, url2)
	assert.Equal(t, expectedURL2.ShortCode, url2.ShortCode)
}

func TestURLUsecase_GetLongURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	usecase := NewURLUsecase(mockRepo)

	shortCode := "testcode"
	longURL := "https://example.com"
	u7, err := uuid.NewV7()
	assert.NoError(t, err)
	expectedURL := &entities.URL{
		ID:          u7.String(),
		ShortCode:   shortCode,
		OriginalURL: longURL,
		CreatedAt:   time.Now(),
	}

	mockRepo.EXPECT().GetURLByShortCode(gomock.Any(), shortCode).Return(expectedURL, nil)
	retrievedLongURL, err := usecase.GetLongURL(context.Background(), shortCode)
	assert.NoError(t, err)
	assert.Equal(t, longURL, retrievedLongURL)

	// Test URL not found
	mockRepo.EXPECT().GetURLByShortCode(gomock.Any(), "nonexistent").Return(nil, errors.New("not found"))
	retrievedLongURL, err = usecase.GetLongURL(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Empty(t, retrievedLongURL)
}

func TestURLUsecase_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	usecase := NewURLUsecase(mockRepo)

	var userID string
	u, err := uuid.NewV7()
	assert.NoError(t, err)
	userID = u.String()

	u7_1, err := uuid.NewV7()
	assert.NoError(t, err)
	u7_2, err := uuid.NewV7()
	assert.NoError(t, err)
	expectedURLs := []*entities.URL{
		{
			ID:          u7_1.String(),
			UserID:      &userID,
			ShortCode:   "code1",
			OriginalURL: "https://example.com/1",
			CreatedAt:   time.Now(),
		},
		{
			ID:          u7_2.String(),
			UserID:      &userID,
			ShortCode:   "code2",
			OriginalURL: "https://example.com/2",
			CreatedAt:   time.Now(),
		},
	}

	mockRepo.EXPECT().GetURLsByUserID(gomock.Any(), userID).Return(expectedURLs, nil)
	urls, err := usecase.GetUserURLs(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedURLs, urls)
}

func TestURLUsecase_DeleteURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	usecase := NewURLUsecase(mockRepo)

	var urlID string
	u, err := uuid.NewV7()
	assert.NoError(t, err)
	urlID = u.String()

	var userID string
	u, err = uuid.NewV7()
	assert.NoError(t, err)
	userID = u.String()

	mockRepo.EXPECT().DeleteURL(gomock.Any(), urlID, userID).Return(nil)
	err = usecase.DeleteURL(context.Background(), urlID, userID)
	assert.NoError(t, err)

	// Test error during deletion
	mockRepo.EXPECT().DeleteURL(gomock.Any(), urlID, userID).Return(errors.New("deletion error"))
	err = usecase.DeleteURL(context.Background(), urlID, userID)
	assert.Error(t, err)
}
