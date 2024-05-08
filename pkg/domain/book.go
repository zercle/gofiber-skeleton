package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type BookUsecase interface {
	GetBook(bookID string) (book models.Book, err error)
	GetBooks(criteria models.Book) (books []models.Book, err error)
	CreateBook(book *models.Book) (err error)
	EditBook(bookID string, book models.Book) (err error)
	DeleteBook(bookID string) (err error)
}

type BookRepository interface {
	DbMigrator() (err error)
	GetBook(bookID string) (book models.Book, err error)
	GetBooks(criteria models.Book) (books []models.Book, err error)
	CreateBook(book *models.Book) (err error)
	EditBook(bookID string, book models.Book) (err error)
	DeleteBook(bookID string) (err error)
}
