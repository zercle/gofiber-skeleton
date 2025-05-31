package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gofiber-skeleton/internal/domain"
	"gofiber-skeleton/internal/handler"
	"gofiber-skeleton/internal/usecase"
	"gofiber-skeleton/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp(svc usecase.BookingService) *fiber.App {
	app := fiber.New()
	handler.RegisterRoutes(app, svc)
	return app
}

func TestCreateBooking_Success(t *testing.T) {
	start := time.Now().UTC().Truncate(time.Second)
	end := start.Add(2 * time.Hour).UTC().Truncate(time.Second)
	payload := map[string]interface{}{
		"user_id":    1,
		"item_id":    2,
		"start_date": start.Format(time.RFC3339),
		"end_date":   end.Format(time.RFC3339),
	}
	body, _ := json.Marshal(payload)

	mockSvc := &mocks.BookingServiceMock{
		CreateBookingFn: func(b *domain.Booking) error {
			return nil
		},
	}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var got domain.Booking
	err := json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), got.UserID)
	assert.Equal(t, uint(2), got.ItemID)
	assert.True(t, got.StartDate.Equal(start))
	assert.True(t, got.EndDate.Equal(end))
}

func TestCreateBooking_InvalidPayload(t *testing.T) {
	mockSvc := &mocks.BookingServiceMock{}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewReader([]byte("not json")))
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetBooking_Success(t *testing.T) {
	existing := &domain.Booking{ID: 10, UserID: 5, ItemID: 6}
	mockSvc := &mocks.BookingServiceMock{
		GetBookingFn: func(id uint) (*domain.Booking, error) {
			return existing, nil
		},
	}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/bookings/10", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var got domain.Booking
	json.NewDecoder(resp.Body).Decode(&got)
	assert.Equal(t, existing.ID, got.ID)
	assert.Equal(t, existing.UserID, got.UserID)
}

func TestGetBooking_NotFound(t *testing.T) {
	mockSvc := &mocks.BookingServiceMock{
		GetBookingFn: func(id uint) (*domain.Booking, error) {
			return nil, usecase.ErrBookingNotFound
		},
	}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/bookings/1", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateBooking_BadID(t *testing.T) {
	mockSvc := &mocks.BookingServiceMock{}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/bookings/0", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateBooking_Success(t *testing.T) {
	start := time.Now().UTC().Truncate(time.Second)
	end := start.Add(time.Hour).UTC().Truncate(time.Second)
	update := map[string]interface{}{
		"user_id":    2,
		"item_id":    3,
		"start_date": start.Format(time.RFC3339),
		"end_date":   end.Format(time.RFC3339),
	}
	body, _ := json.Marshal(update)

	mockSvc := &mocks.BookingServiceMock{
		UpdateBookingFn: func(b *domain.Booking) error {
			return nil
		},
	}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/bookings/5", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var got domain.Booking
	json.NewDecoder(resp.Body).Decode(&got)
	assert.Equal(t, uint(2), got.UserID)
	assert.Equal(t, uint(3), got.ItemID)
}

func TestDeleteBooking_Success(t *testing.T) {
	mockSvc := &mocks.BookingServiceMock{
		DeleteBookingFn: func(id uint) error {
			return nil
		},
	}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/bookings/7", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDeleteBooking_BadID(t *testing.T) {
	mockSvc := &mocks.BookingServiceMock{}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/bookings/0", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestListBookingsByUser_Success(t *testing.T) {
	bookings := []domain.Booking{{ID: 1, UserID: 8}, {ID: 2, UserID: 8}}
	mockSvc := &mocks.BookingServiceMock{
		ListBookingsByUserFn: func(userID uint) ([]domain.Booking, error) {
			return bookings, nil
		},
	}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/users/8/bookings", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var got []domain.Booking
	json.NewDecoder(resp.Body).Decode(&got)
	assert.Equal(t, bookings, got)
}

func TestListBookingsByUser_BadID(t *testing.T) {
	mockSvc := &mocks.BookingServiceMock{}
	app := setupApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/users/0/bookings", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
