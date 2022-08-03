package entities

import (
	"time"

	"gorm.io/gorm"
)

// Book Constructs your Book model under entities.
type Book struct {
	Id          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:127;index"`
	Author      string         `json:"author" gorm:"size:127;index"`
	Description string         `json:"description" gorm:""`
	CreatedAt   time.Time      `json:"createdAt" gorm:"autoCreateTime;index"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"autoUpdateTime;index"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// DeleteRequest struct is used to parse Delete Reqeusts for Books
type DeleteRequest struct {
	Id string `json:"id"`
}
