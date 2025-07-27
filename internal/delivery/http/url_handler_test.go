package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gofiber-skeleton/internal/entities"
	"gofiber-skeleton/internal/usecases/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_CreateShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_ = context.Background() // Dummy usage to prevent "context imported and not used" error

	mockURLUseCase := mocks.NewMockURLUseCase(ctrl)
	handler := NewURLHandler(mockURLUseCase)

	var currentUserID uuid.UUID // Declare userID in the outer scope

	app := fiber.New()
	app.Post("/urls", func(c *fiber.Ctx) error {
		// Mock JWT middleware
		c.Locals("user", &jwt.Token{
			Claims: jwt.MapClaims{
				"sub": currentUserID.String(), // Use the captured userID
			},
		})
		return handler.CreateShortURL(c)
	})

	tests := []struct {
		name           string
		requestBody    map[string]string
		mockSetup      func(userID uuid.UUID)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful short URL creation",
			requestBody: map[string]string{
				"original_url": "https://example.com",
			},
			mockSetup: func(userID uuid.UUID) {
				mockURLUseCase.EXPECT().CreateShortURL(gomock.Any(), "https://example.com", userID, "").Return(&entities.URL{
					OriginalURL: "https://example.com",
					ShortCode:   "shortcode123",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":"success","message":"Short URL created successfully","data":{"short_code":"shortcode123"}}`,
		},
		{
			name: "invalid request body",
			requestBody: map[string]string{
				"invalid_field": "value",
			},
			mockSetup:      func(userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","message":"Original URL cannot be empty"}`,
		},
		{
			name: "error from use case",
			requestBody: map[string]string{
				"original_url": "https://example.com",
			},
			mockSetup: func(userID uuid.UUID) {
				mockURLUseCase.EXPECT().CreateShortURL(gomock.Any(), "https://example.com", userID, "").Return(nil, errors.New("failed to create short URL"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"error","message":"Failed to create short URL"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentUserID = uuid.New() // Assign new userID for each test case
			tt.mockSetup(currentUserID)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/urls", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generateMockToken(currentUserID)))

			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBody, string(respBody))
		})
	}
}

func TestURLHandler_GetOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLUseCase := mocks.NewMockURLUseCase(ctrl)
	handler := NewURLHandler(mockURLUseCase)

	app := fiber.New()
	app.Get("/:shortCode", handler.Redirect)

	tests := []struct {
		name           string
		shortCode      string
		mockSetup      func()
		expectedStatus int
		expectedHeader map[string]string
	}{
		{
			name:      "successful redirection",
			shortCode: "testcode",
			mockSetup: func() {
				mockURLUseCase.EXPECT().GetOriginalURL(gomock.Any(), "testcode").Return("https://example.com/original", nil)
			},
			expectedStatus: http.StatusMovedPermanently,
			expectedHeader: map[string]string{"Location": "https://example.com/original"},
		},
		{
			name:      "short code not found",
			shortCode: "notfound",
			mockSetup: func() {
				mockURLUseCase.EXPECT().GetOriginalURL(gomock.Any(), "notfound").Return("", errors.New("URL not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedHeader: map[string]string{"Content-Type": "application/json"}, // Default Fiber error response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.shortCode), nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			for key, value := range tt.expectedHeader {
				assert.Equal(t, value, resp.Header.Get(key))
			}
		})
	}
}


func TestURLHandler_GenerateQRCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockURLUseCase := mocks.NewMockURLUseCase(ctrl)
	handler := NewURLHandler(mockURLUseCase)

	app := fiber.New()
	app.Get("/:shortCode/qr", handler.GetQRCode)

	tests := []struct {
		name           string
		shortCode      string
		mockSetup      func()
		expectedStatus int
		expectedHeader map[string]string
		expectedBody   []byte
	}{
		{
			name:      "successful QR code generation",
			shortCode: "testcode",
			mockSetup: func() {
				mockURLUseCase.EXPECT().GenerateQRCode(gomock.Any(), "testcode").Return([]byte("mock_qr_code_image"), nil)
			},
			expectedStatus: http.StatusOK,
			expectedHeader: map[string]string{"Content-Type": "image/png"},
			expectedBody:   []byte("mock_qr_code_image"),
		},
		{
			name:      "error generating QR code",
			shortCode: "errorcode",
			mockSetup: func() {
				mockURLUseCase.EXPECT().GenerateQRCode(gomock.Any(), "errorcode").Return(nil, errors.New("failed to generate QR code"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
			expectedBody:   []byte(`{"status":"error","message":"Failed to generate QR code"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/qr", tt.shortCode), nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			for key, value := range tt.expectedHeader {
				assert.Equal(t, value, resp.Header.Get(key))
			}
			respBody, _ := io.ReadAll(resp.Body)
			if tt.expectedHeader["Content-Type"] == "image/png" {
				assert.Equal(t, tt.expectedBody, respBody)
			} else {
				assert.JSONEq(t, string(tt.expectedBody), string(respBody))
			}
		})
	}
}

// Helper function to generate a mock JWT token for testing
func generateMockToken(userID uuid.UUID) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))
	return tokenString
}