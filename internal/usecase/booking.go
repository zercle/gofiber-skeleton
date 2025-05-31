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

	 // BookingService defines use-case methods for bookings
	 type BookingService interface {
	 	CreateBooking(b *domain.Booking) error
	 	GetBooking(id uint) (*domain.Booking, error)
	 	UpdateBooking(b *domain.Booking) error
	 	DeleteBooking(id uint) error
	 	ListBookingsByUser(userID uint) ([]domain.Booking, error)
	 }

type BookingUsecase struct {
	repo domain.BookingRepository
}

func NewBookingUsecase(r domain.BookingRepository) BookingService {
	return &BookingUsecase{repo: r}
}

func (uc *BookingUsecase) CreateBooking(b *domain.Booking) error {
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

func (uc *BookingUsecase) GetBooking(id uint) (*domain.Booking, error) {
	if id == 0 {
		return nil, ErrInvalidBooking
	}
	booking, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, ErrBookingNotFound
	}
	return booking, nil
}

func (uc *BookingUsecase) UpdateBooking(b *domain.Booking) error {
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

func (uc *BookingUsecase) DeleteBooking(id uint) error {
	if id == 0 {
		return ErrInvalidBooking
	}
	return uc.repo.Delete(id)
}

func (uc *BookingUsecase) ListBookingsByUser(userID uint) ([]domain.Booking, error) {
	if userID == 0 {
		return nil, ErrInvalidBooking
	}
	return uc.repo.ListByUser(userID)
}