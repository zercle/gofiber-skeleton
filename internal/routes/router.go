package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/books"
	"gorm.io/gorm"
)

type RouterResources struct {
	DbConn *gorm.DB
}

func NewRouterResources(dbConn *gorm.DB) (resources *RouterResources) {
	resources = &RouterResources{
		DbConn: dbConn,
	}
	return
}

// SetupRoutes is the Router for GoFiber App
func (r *RouterResources) SetupRoutes(app *fiber.App) {

	app.Get("/", handlers.Index())

	groupApiV1 := app.Group("/api/v:version?", apiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}

	// App Repository
	bookRepository := books.NewBookRepository(r.DbConn)

	// App Services
	bookService := books.NewBookService(bookRepository)

	// App Routes
	books.NewBookHandler(app.Group("/api/v1/book"), bookService)
}
