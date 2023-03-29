package books

import (
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type bookUsecase struct {
	bookRepository domain.BookRepository
}

func NewBookUsecase(r domain.BookRepository) domain.BookUsecase {
	return &bookUsecase{
		bookRepository: r,
	}
}

func (u *bookUsecase) DbMigrator() (err error) {
	err = u.bookRepository.DbMigrator()
	return
}

func (u *bookUsecase) CreateBook(book *models.Book) (err error) {
	if len(book.Title) == 0 {
		err = helpers.NewError(fiber.StatusBadRequest, helpers.WhereAmI(), helpers.WhereAmI(), "need: title")
		return
	}
	if len(book.Author) == 0 {
		err = helpers.NewError(fiber.StatusBadRequest, helpers.WhereAmI(), helpers.WhereAmI(), "need: author")
		return
	}
	return u.bookRepository.CreateBook(book)
}

func (u *bookUsecase) EditBook(bookId uint, book models.Book) (err error) {
	return u.bookRepository.EditBook(bookId, book)
}

func (u *bookUsecase) DeleteBook(bookId uint) (err error) {
	return u.bookRepository.DeleteBook(bookId)
}

func (u *bookUsecase) GetBook(bookId uint) (book models.Book, err error) {
	return u.bookRepository.GetBook(bookId)
}

func (u *bookUsecase) GetBooks(criteria models.Book) (books []models.Book, err error) {
	return u.bookRepository.GetBooks(criteria)
}
