package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/books"
	"github.com/zercle/gofiber-skelton/pkg/users"
)

// SetupRoutes is the Router for GoFiber App
func (s *Server) SetupRoutes(app *fiber.App) {

	// Prepare a static middleware to serve the built React files.
	app.Static("/", "./web/build")

	// API routes group
	groupApiV1 := app.Group("/api/v:version?", handlers.ApiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}

	// App Repository
	bookRepository := books.InitBookRepository(s.MainDbConn)
	userRepository := users.InitUserRepository(s.MainDbConn)

	// App Services
	bookUsecase := books.InitBookUsecase(bookRepository)
	userUsecase := users.InitUserUsecase(userRepository)

	// App Routes
	books.NewBookHandler(app.Group("/api/v1/books"), bookUsecase)
	users.InitUserHandler(app.Group("/api/v1/users"), userUsecase)

	// Prepare a fallback route to always serve the 'index.html', had there not be any matching routes.
	app.Static("*", "./web/build/index.html")
}
