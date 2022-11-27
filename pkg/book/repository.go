package book

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type bookReposiroty struct {
	DbConn *gorm.DB
}

func NewBookRepository(dbConn *gorm.DB) BookRepository {
	return &bookReposiroty{
		DbConn: dbConn,
	}
}

func (r *bookReposiroty) CreateBook(book *Book) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&Book{})

	if err = dbTx.Create(book).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error
	return
}

func (r *bookReposiroty) UpdateBook(bookId uint, book *Book) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&Book{})
	dbTx = dbTx.Where(Book{Id: bookId})

	if err = dbTx.Updates(book).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *bookReposiroty) DeleteBook(bookId uint) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&Book{})
	dbTx = dbTx.Where(Book{Id: bookId})

	if err = dbTx.Delete(&Book{}).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *bookReposiroty) GetBook(bookId uint) (book *Book, err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Model(&Book{})
	dbTx = dbTx.Where(Book{Id: bookId})
	err = dbTx.Take(book).Error

	return
}

func (r *bookReposiroty) GetBooks(title, author string) (books []Book, err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Model(&Book{})

	if len(title) != 0 {
		dbTx = dbTx.Where("title LIKE ?", "%"+title+"%")
	}

	if len(author) != 0 {
		dbTx = dbTx.Where("author LIKE ?", "%"+author+"%")
	}

	err = dbTx.Find(&books).Error

	return
}
