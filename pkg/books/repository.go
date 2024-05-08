package books

import (
	"runtime"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type bookReposiroty struct {
	mainDbConn *gorm.DB
}

func NewBookRepository(mainDbConn *gorm.DB) domain.BookRepository {
	return &bookReposiroty{
		mainDbConn: mainDbConn,
	}
}

func (r *bookReposiroty) DbMigrator() (err error) {
	err = r.mainDbConn.AutoMigrate(&models.Book{})
	return
}

func (r *bookReposiroty) CreateBook(book *models.Book) (err error) {
	if r.mainDbConn == nil {
		err = helpers.NewError(fiber.StatusServiceUnavailable, helpers.WhereAmI(), "Database server has gone away")
		return
	}

	dbTx := r.mainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.Book{})

	if err = dbTx.Create(book).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error
	return
}

func (r *bookReposiroty) EditBook(bookID string, book models.Book) (err error) {
	if r.mainDbConn == nil {
		err = helpers.NewError(fiber.StatusServiceUnavailable, helpers.WhereAmI(), "Database server has gone away")
		return
	}

	dbTx := r.mainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.Book{})
	dbTx = dbTx.Where(models.Book{ID: bookID})

	if err = dbTx.Updates(book).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *bookReposiroty) DeleteBook(bookID string) (err error) {
	if r.mainDbConn == nil {
		err = helpers.NewError(fiber.StatusServiceUnavailable, helpers.WhereAmI(), "Database server has gone away")
		return
	}

	dbTx := r.mainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.Book{})
	dbTx = dbTx.Where(models.Book{ID: bookID})

	if err = dbTx.Delete(&models.Book{}).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *bookReposiroty) GetBook(bookID string) (book models.Book, err error) {
	if r.mainDbConn == nil {
		err = helpers.NewError(fiber.StatusServiceUnavailable, helpers.WhereAmI(), "Database server has gone away")
		return
	}

	dbTx := r.mainDbConn.Model(&models.Book{})
	dbTx = dbTx.Where(models.Book{ID: bookID})
	err = dbTx.Take(&book).Error

	return
}

func (r *bookReposiroty) GetBooks(criteria models.Book) (books []models.Book, err error) {
	if r.mainDbConn == nil {
		err = helpers.NewError(fiber.StatusServiceUnavailable, helpers.WhereAmI(), "Database server has gone away")
		return
	}

	dbTx := r.mainDbConn.Model(&models.Book{})

	if len(criteria.ID) != 0 {
		dbTx = dbTx.Where(models.Book{ID: criteria.ID})
	} else {
		if len(criteria.Title) != 0 {
			dbTx = dbTx.Where("title LIKE ?", "%"+criteria.Title+"%")
		}

		if len(criteria.Author) != 0 {
			dbTx = dbTx.Where("author LIKE ?", "%"+criteria.Author+"%")
		}
	}

	err = dbTx.Find(&books).Error

	return
}

func (r *bookReposiroty) ImportBooks(books []models.Book) (errs []error) {

	errCh := make(chan error, (runtime.NumCPU()/2)+1)
	defer close(errCh)

	for _, book := range books {
		go r.importBookChannel(book, errCh)
		if err := <-errCh; err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (r *bookReposiroty) importBookChannel(book models.Book, errCh chan error) {
	errCh <- r.CreateBook(&book)
}
