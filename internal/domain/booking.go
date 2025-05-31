package domain

import "time"

type Booking struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	ItemID    uint      `gorm:"not null"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BookingRepository interface {
	// Create persists a new booking record.
	Create(booking *Booking) error
	// GetByID retrieves a booking by its ID.
	GetByID(id uint) (*Booking, error)
	// Update modifies an existing booking.
	Update(booking *Booking) error
	// Delete removes a booking by its ID.
	Delete(id uint) error
	// ListByUser fetches all bookings for a given user.
	ListByUser(userID uint) ([]Booking, error)
}