package domain

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"size:255;not null;"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:255;not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate is a GORM hook that is called before a new record is created.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = id
	return
}
