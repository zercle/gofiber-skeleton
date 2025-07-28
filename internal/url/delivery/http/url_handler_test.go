package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gofiber-skeleton/internal/url"
	urlHandler "gofiber-skeleton/internal/url/delivery/http"
	mock_usecase "gofiber-skeleton/internal/url/usecase/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_CreateURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLUseCase := mock_usecase.NewMockURLUseCase(ctrl)
	handler := urlHandler.NewHTTPURLHandler(mockURLUseCase)

	app := fiber.New()
	app.Post("/api/urls", handler.CreateURL)

	// Test case 1: Successful URL creation
	originalURL := "https://example.com"
	shortCode := "abc123def"
	userID := uuid.New()
	expectedURL := &url.ModelURL{
		ID:          uuid.New(),
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	mockURLUseCase.EXPECT().CreateURL(gomock.Any(), originalURL, userID.String()).Return(expectedURL, nil).Times(1)

	createURLPayload := map[string]string{
		"original_url": originalURL,
		"user_id":      userID.String(),
	}
	body, _ := json.Marshal(createURLPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, shortCode, data["short_code"])

	// Test case 2: Invalid request body
	req = httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test case 3: Use case error
	useCaseError := errors.New("failed to process URL")
	mockURLUseCase.EXPECT().CreateURL(gomock.Any(), originalURL, userID.String()).Return(nil, useCaseError).Times(1)

	body, _ = json.Marshal(createURLPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestURLHandler_GetShortenedURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLUseCase := mock_usecase.NewMockURLUseCase(ctrl)
	handler := urlHandler.NewHTTPURLHandler(mockURLUseCase)

	app := fiber.New()
	app.Get("/:shortCode", handler.GetOriginalURL)

	// Test case 1: Successful redirection
	shortCode := "validcode"
	originalURL := "https://redirect.com"
	mockURLUseCase.EXPECT().GetOriginalURL(gomock.Any(), shortCode).Return(originalURL, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
	assert.Equal(t, originalURL, resp.Header.Get("Location"))

	// Test case 2: Use case error (URL not found/expired)
	useCaseError := errors.New("URL not found or expired")
	mockURLUseCase.EXPECT().GetOriginalURL(gomock.Any(), "invalidcode").Return("", useCaseError).Times(1)

	req = httptest.NewRequest(http.MethodGet, "/invalidcode", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
