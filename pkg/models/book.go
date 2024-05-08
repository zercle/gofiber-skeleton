package models

import (
	"database/sql"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type BooksResponse struct {
	Books []Book `json:"books"`
}

// Book Constructs your Book model under entities.
type Book struct {
	CreatedAt   sql.NullTime   `json:"createdAt,omitempty" gorm:"autoCreateTime;index"`
	UpdatedAt   sql.NullTime   `json:"updatedAt,omitempty" gorm:"autoUpdateTime;index"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
	Title       string         `json:"title" gorm:"size:127;index"`
	Author      string         `json:"author" gorm:"size:127;index"`
	Description string         `json:"description" gorm:""`
	ID          string         `json:"id" gorm:"size:32;primaryKey"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	if len(b.ID) == 0 {
		b.ID = ulid.Make().String()
	}
	return
}
