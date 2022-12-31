package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type BooksResponse struct {
	Books []Book `json:"books"`
}

// Book Constructs your Book model under entities.
type Book struct {
	Id          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:127;index"`
	Author      string         `json:"author" gorm:"size:127;index"`
	Description string         `json:"description" gorm:""`
	CreatedAt   sql.NullTime   `json:"createdAt,omitempty" gorm:"autoCreateTime;index"`
	UpdatedAt   sql.NullTime   `json:"updatedAt,omitempty" gorm:"autoUpdateTime;index"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}
