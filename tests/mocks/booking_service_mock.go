package mocks

import (
	"gofiber-skeleton/internal/domain"
)

type BookingServiceMock struct {
	CreateBookingFn       func(b *domain.Booking) error
	GetBookingFn          func(id uint) (*domain.Booking, error)
	UpdateBookingFn       func(b *domain.Booking) error
	DeleteBookingFn       func(id uint) error
	ListBookingsByUserFn  func(userID uint) ([]domain.Booking, error)
}

func (m *BookingServiceMock) CreateBooking(b *domain.Booking) error {
	if m.CreateBookingFn != nil {
		return m.CreateBookingFn(b)
	}
	return nil
}

func (m *BookingServiceMock) GetBooking(id uint) (*domain.Booking, error) {
	if m.GetBookingFn != nil {
		return m.GetBookingFn(id)
	}
	return &domain.Booking{}, nil
}

func (m *BookingServiceMock) UpdateBooking(b *domain.Booking) error {
	if m.UpdateBookingFn != nil {
		return m.UpdateBookingFn(b)
	}
	return nil
}

func (m *BookingServiceMock) DeleteBooking(id uint) error {
	if m.DeleteBookingFn != nil {
		return m.DeleteBookingFn(id)
	}
	return nil
}

func (m *BookingServiceMock) ListBookingsByUser(userID uint) ([]domain.Booking, error) {
	if m.ListBookingsByUserFn != nil {
		return m.ListBookingsByUserFn(userID)
	}
	return nil, nil
}