package books

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/pkg/models"
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

func (r *bookReposiroty) dbMigrator() (err error) {
	err = r.DbConn.AutoMigrate(&models.Book{})
	return
}

func (r *bookReposiroty) CreateBook(book *models.Book) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.Book{})

	if err = dbTx.Create(book).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error
	return
}

func (r *bookReposiroty) UpdateBook(bookId uint, book *models.Book) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.Book{})
	dbTx = dbTx.Where(models.Book{Id: bookId})

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

	dbTx = dbTx.Model(&models.Book{})
	dbTx = dbTx.Where(models.Book{Id: bookId})

	if err = dbTx.Delete(&models.Book{}).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *bookReposiroty) GetBook(bookId uint) (book *models.Book, err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Model(&models.Book{})
	dbTx = dbTx.Where(models.Book{Id: bookId})
	err = dbTx.Take(book).Error

	return
}

func (r *bookReposiroty) GetBooks(title, author string) (books []models.Book, err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Model(&models.Book{})

	if len(title) != 0 {
		dbTx = dbTx.Where("title LIKE ?", "%"+title+"%")
	}

	if len(author) != 0 {
		dbTx = dbTx.Where("author LIKE ?", "%"+author+"%")
	}

	err = dbTx.Find(&books).Error

	return
}
