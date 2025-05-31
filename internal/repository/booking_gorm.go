package repository

import (
	"gofiber-skeleton/internal/domain"
	"gorm.io/gorm"
)

type bookingGormRepo struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) domain.BookingRepository {
	return &bookingGormRepo{db: db}
}

func (r *bookingGormRepo) Create(b *domain.Booking) error {
	return r.db.Create(b).Error
}

func (r *bookingGormRepo) GetByID(id uint) (*domain.Booking, error) {
	var booking domain.Booking
	if err := r.db.First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingGormRepo) Update(b *domain.Booking) error {
	return r.db.Save(b).Error
}

func (r *bookingGormRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Booking{}, id).Error
}

func (r *bookingGormRepo) ListByUser(userID uint) ([]domain.Booking, error) {
	var bookings []domain.Booking
	if err := r.db.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}