package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/handlers"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
)

func SetupPostRoutes(app *fiber.App, postHandler *handlers.PostHandler, cfg *config.Config) {
	api := app.Group("/api/v1")
	posts := api.Group("/posts")

	posts.Get("/", postHandler.ListPosts)
	posts.Get("/:id", postHandler.GetPost)

	authMiddleware := middleware.NewAuthMiddleware(cfg)
	
	posts.Use(authMiddleware)
	posts.Post("/", postHandler.CreatePost)
	posts.Put("/:id", postHandler.UpdatePost)
	posts.Delete("/:id", postHandler.DeletePost)
	posts.Get("/me", postHandler.ListUserPosts)
	posts.Post("/:id/publish", postHandler.PublishPost)
	posts.Post("/:id/unpublish", postHandler.UnpublishPost)
}