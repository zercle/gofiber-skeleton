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

func (u *bookUsecase) EditBook(bookID uint, book models.Book) (err error) {
	return u.bookRepository.EditBook(bookID, book)
}

func (u *bookUsecase) DeleteBook(bookID uint) (err error) {
	return u.bookRepository.DeleteBook(bookID)
}

func (u *bookUsecase) GetBook(bookID uint) (book models.Book, err error) {
	return u.bookRepository.GetBook(bookID)
}

func (u *bookUsecase) GetBooks(criteria models.Book) (books []models.Book, err error) {
	return u.bookRepository.GetBooks(criteria)
}
