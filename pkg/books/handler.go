package books

import (
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type BookHandler struct {
	bookService BookService
}

func NewBookHandler(bookRoute fiber.Router, bs BookService) {
	handler := &BookHandler{
		bookService: bs,
	}

	bookRoute.Get("/:bookId?", handler.getBooks())
}

func (h *BookHandler) getBooks() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}

		bookId, _ := c.ParamsInt("bookId", 0)
		if bookId != 0 {
			book, err := h.bookService.GetBook(uint(bookId))
			if err != nil {
				responseForm.Errors = append(responseForm.Errors, &helpers.ResposeError{
					Code:    fiber.StatusServiceUnavailable,
					Message: err.Error(),
				})
			}
			responseForm.Data = BooksResponse{
				Books: &[]Book{*book},
			}
		} else {
			title := c.FormValue("title")
			author := c.FormValue("author")
			books, err := h.bookService.GetBooks(title, author)
			if err != nil {
				responseForm.Errors = append(responseForm.Errors, &helpers.ResposeError{
					Code:    fiber.StatusServiceUnavailable,
					Message: err.Error(),
				})
			}
			responseForm.Data = BooksResponse{
				Books: books,
			}
		}

		if err == nil {
			responseForm.Success = true
		}
		return c.JSON(responseForm)
	}
}

func (h *BookHandler) createBook() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}

		book := new(Book)

		if err = c.BodyParser(book); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			responseForm.Errors = append(responseForm.Errors, &helpers.ResposeError{
				Code:    fiber.StatusUnprocessableEntity,
				Message: err.Error(),
			})
		}

		if err != nil {
			return c.JSON(responseForm)
		}

		bookId, _ := c.ParamsInt("bookId", 0)
		if bookId != 0 {
			book, err := h.bookService.GetBook(uint(bookId))
			if err != nil {
				c.Status(fiber.StatusServiceUnavailable)
				responseForm.Errors = append(responseForm.Errors, &helpers.ResposeError{
					Code:    fiber.StatusServiceUnavailable,
					Message: err.Error(),
				})
			}
			responseForm.Result = BooksResponse{
				Books: &[]Book{*book},
			}
		} else {
			title := c.FormValue("title")
			author := c.FormValue("author")
			books, err := h.bookService.GetBooks(title, author)
			if err != nil {
				c.Status(fiber.StatusServiceUnavailable)
				responseForm.Errors = append(responseForm.Errors, &helpers.ResposeError{
					Code:    fiber.StatusServiceUnavailable,
					Message: err.Error(),
				})
			}
			responseForm.Result = BooksResponse{
				Books: books,
			}
		}

		if err == nil {
			responseForm.Success = true
		}
		return c.JSON(responseForm)
	}
}
