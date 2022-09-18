package books

import (
	"time"

	"gorm.io/gorm"
)

type BooksResponse struct {
	Books *[]Book `json:"books"`
}

// Book Constructs your Book model under entities.
type Book struct {
	Id          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:127;index"`
	Author      string         `json:"author" gorm:"size:127;index"`
	Description string         `json:"description" gorm:""`
	CreatedAt   *time.Time     `json:"createdAt,omitempty" gorm:"autoCreateTime;index"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty" gorm:"autoUpdateTime;index"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type BookRepository interface {
	GetBooks(title, author string) (*[]Book, error)
	GetBook(bookID uint) (*Book, error)
	CreateBook(book *Book) error
	UpdateBook(bookID uint, book *Book) error
	DeleteBook(bookID uint) error
}

type BookService interface {
	GetBooks(title, author string) (*[]Book, error)
	GetBook(bookID uint) (*Book, error)
	CreateBook(book *Book) error
	UpdateBook(bookID uint, book *Book) error
	DeleteBook(bookID uint) error
}
