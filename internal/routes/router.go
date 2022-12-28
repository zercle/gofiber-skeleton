package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/books"
)

// SetupRoutes is the Router for GoFiber App
func (r *RouterResources) SetupRoutes(app *fiber.App) {

	// Prepare a static middleware to serve the built React files.
	app.Static("/", "./web/build")

	// API routes group
	groupApiV1 := app.Group("/api/v:version?", apiLimiter)
	{
		groupApiV1.Get("/", r.Index())
	}

	// App Repository
	bookRepository := books.NewBookRepository(r.MainDbConn)

	// App Services
	bookService := books.NewBookService(bookRepository)

	// App Routes
	books.NewBookHandler(app.Group("/api/v1/books"), bookService)

	// Prepare a fallback route to always serve the 'index.html', had there not be any matching routes.
	app.Static("*", "./web/build/index.html")
}
