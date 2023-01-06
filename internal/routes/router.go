package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/books"
	"github.com/zercle/gofiber-skelton/internal/users"
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
	bookRepository := books.InitBookRepository(r.Resources)
	userRepository := users.InitUserRepository(r.Resources)

	// App Services
	bookUsecase := books.InitBookUsecase(bookRepository)
	userUsecase := users.InitUserUsecase(userRepository)

	// App Routes
	books.NewBookHandler(app.Group("/api/v1/books"), bookUsecase)
	users.NewUserHandler(app.Group("/api/v1/users"), userUsecase)

	// Prepare a fallback route to always serve the 'index.html', had there not be any matching routes.
	app.Static("*", "./web/build/index.html")
}
