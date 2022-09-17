package books

import "github.com/gofiber/fiber/v2"

type BookHandler struct {
	bookService BookService
}

func NewBookHandler(bookRoute fiber.Router, bs BookService) {
	handler := &BookHandler{
		bookService: bs,
	}

	bookRoute.Get("", handler.getBook())
}

func (h *BookHandler) getBook() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		return
	}
}
