package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type BookUsecase interface {
	GetBook(bookId uint) (book models.Book, err error)
	GetBooks(criteria models.Book) (books []models.Book, err error)
	CreateBook(book *models.Book) (err error)
	EditBook(bookId uint, book models.Book) (err error)
	DeleteBook(bookId uint) (err error)
}

type BookRepository interface {
	DbMigrator() (err error)
	GetBook(bookId uint) (book models.Book, err error)
	GetBooks(criteria models.Book) (books []models.Book, err error)
	CreateBook(book *models.Book) (err error)
	EditBook(bookId uint, book models.Book) (err error)
	DeleteBook(bookId uint) (err error)
}
