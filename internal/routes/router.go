package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/pkg/book"
)

// SetupRoutes is the Router for GoFiber App
func (r *RouterResources) SetupRoutes(app *fiber.App) {

	app.Get("/", r.Index())

	groupApiV1 := app.Group("/api/v:version?", apiLimiter)
	{
		groupApiV1.Get("/", r.Index())
	}

	// App Repository
	bookRepository := book.NewBookRepository(r.MainDbConn)

	// App Services
	bookService := book.NewBookService(bookRepository)

	// App Routes
	book.NewBookHandler(app.Group("/api/v1/books"), bookService)
}
