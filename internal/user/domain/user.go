package domain

import (
	"time"

	"gofiber-skeleton/pkg/types"

	"gorm.io/gorm"
)

// User represents the user model.
type User struct {
	ID        types.UUIDv7   `json:"id" gorm:"type:uuid;primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}