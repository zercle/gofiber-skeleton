package books

import (
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type BookHandler struct {
	bookUsecase domain.BookUsecase
}

func NewBookHandler(bookRoute fiber.Router, bookUsecase domain.BookUsecase) {

	handler := &BookHandler{
		bookUsecase: bookUsecase,
	}

	bookRoute.Get("/:bookID?", handler.getBooks())
}

func (h *BookHandler) getBooks() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}

		bookID, _ := c.ParamsInt("bookID", 0)
		if bookID != 0 {
			book, err := h.bookUsecase.GetBook(uint(bookID))
			if err != nil {
				responseForm.Errors = append(responseForm.Errors, helpers.ResponseError{
					Code:    fiber.StatusServiceUnavailable,
					Message: err.Error(),
				})
			}
			responseForm.Result = models.BooksResponse{
				Books: []models.Book{book},
			}
		} else {
			criteria := models.Book{}
			criteria.Title = c.FormValue("title")
			criteria.Author = c.FormValue("author")
			books, err := h.bookUsecase.GetBooks(criteria)
			if err != nil {
				responseForm.Errors = append(responseForm.Errors, helpers.ResponseError{
					Code:    fiber.StatusServiceUnavailable,
					Message: err.Error(),
				})
			}
			responseForm.Result = models.BooksResponse{
				Books: books,
			}
		}

		if err == nil {
			responseForm.Success = true
		}
		return c.JSON(responseForm)
	}
}
