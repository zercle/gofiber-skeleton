package infrastructure

import (
	"log"

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
	bookRepository := books.NewBookRepository(s.MainDbConn)
	userRepository := users.NewUserRepository(s.MainDbConn)

	// auto migrate DB only on main process
	if !fiber.IsChild() {
		if migrateErr := bookRepository.DbMigrator(); migrateErr != nil {
			log.Panicf("error while migrate book DB:\n %+v", migrateErr)
		}
	}

	// App Services
	bookUsecase := books.NewBookUsecase(bookRepository)
	userUsecase := users.NewUserUsecase(userRepository)

	// App Routes
	books.NewBookHandler(app.Group("/api/v1/books"), bookUsecase)
	users.NewUserHandler(app.Group("/api/v1/users"), userUsecase)

	// Prepare a fallback route to always serve the 'index.html', had there not be any matching routes.
	app.Static("*", "./web/build/index.html")
}
