package usecase

import (
	"errors"
	"gofiber-skeleton/internal/domain"
)

var (
	ErrInvalidBooking      = errors.New("invalid booking data")
	ErrBookingNotFound     = errors.New("booking not found")
	ErrBookingDateConflict = errors.New("start date must be before end date")
	)

	 // BookingUsecase defines use-case methods for bookings
	 type BookingUsecase interface {
	 	CreateBooking(b *domain.Booking) error
	 	GetBooking(id uint) (*domain.Booking, error)
	 	UpdateBooking(b *domain.Booking) error
	 	DeleteBooking(id uint) error
	 	ListBookingsByUser(userID uint) ([]domain.Booking, error)
	 }

type bookingUsecase struct {
	repo domain.BookingRepository
}

func NewBookingUsecase(r domain.BookingRepository) BookingUsecase {
	return &bookingUsecase{repo: r}
}

func (uc *bookingUsecase) CreateBooking(b *domain.Booking) error {
	if b == nil {
		return ErrInvalidBooking
	}
	if b.UserID == 0 || b.ItemID == 0 {
		return ErrInvalidBooking
	}
	if !b.StartDate.Before(b.EndDate) {
		return ErrBookingDateConflict
	}
	return uc.repo.Create(b)
}

func (uc *bookingUsecase) GetBooking(id uint) (*domain.Booking, error) {
	if id == 0 {
		return nil, ErrInvalidBooking
	}
	booking, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, ErrBookingNotFound
	}
	return booking, nil
}

func (uc *bookingUsecase) UpdateBooking(b *domain.Booking) error {
	if b == nil || b.ID == 0 {
		return ErrInvalidBooking
	}
	if b.UserID == 0 || b.ItemID == 0 {
		return ErrInvalidBooking
	}
	if !b.StartDate.Before(b.EndDate) {
		return ErrBookingDateConflict
	}
	return uc.repo.Update(b)
}

func (uc *bookingUsecase) DeleteBooking(id uint) error {
	if id == 0 {
		return ErrInvalidBooking
	}
	return uc.repo.Delete(id)
}

func (uc *bookingUsecase) ListBookingsByUser(userID uint) ([]domain.Booking, error) {
	if userID == 0 {
		return nil, ErrInvalidBooking
	}
	return uc.repo.ListByUser(userID)
}