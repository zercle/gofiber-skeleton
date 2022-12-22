package books

import (
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type BookRepository interface {
	dbMigrator() (err error)

	GetBooks(title, author string) ([]models.Book, error)
	GetBook(bookID uint) (*models.Book, error)
	CreateBook(book *models.Book) error
	UpdateBook(bookID uint, book *models.Book) error
	DeleteBook(bookID uint) error
}

type BookService interface {
	dbMigrator() (err error)

	GetBooks(title, author string) ([]models.Book, error)
	GetBook(bookID uint) (*models.Book, error)
	CreateBook(book *models.Book) error
	UpdateBook(bookID uint, book *models.Book) error
	DeleteBook(bookID uint) error
}
