package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type BookUsecase interface {
	GetBook(bookID uint) (book models.Book, err error)
	GetBooks(criteria models.Book) (books []models.Book, err error)
	CreateBook(book *models.Book) (err error)
	EditBook(bookID uint, book models.Book) (err error)
	DeleteBook(bookID uint) (err error)
}

type BookRepository interface {
	DbMigrator() (err error)
	GetBook(bookID uint) (book models.Book, err error)
	GetBooks(criteria models.Book) (books []models.Book, err error)
	CreateBook(book *models.Book) (err error)
	EditBook(bookID uint, book models.Book) (err error)
	DeleteBook(bookID uint) (err error)
}
