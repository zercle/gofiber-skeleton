package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/api"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/repositories"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/usecases"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
)

// RegisterRoutes registers all posts module routes.
func RegisterRoutes(router fiber.Router, db *database.Database, authMiddleware fiber.Handler) {
	// Initialize repository, usecase, and handler
	postRepo := repositories.NewPostRepository(db)
	postUseCase := usecases.NewPostUseCase(postRepo)
	postHandler := api.NewPostHandler(postUseCase)

	posts := router.Group("/posts")

	// Public routes
	posts.Get("/", postHandler.ListPosts)               // List all posts (can filter by published)
	posts.Get("/:id", postHandler.GetPost)              // Get post by ID
	posts.Get("/slug/:slug", postHandler.GetPostBySlug) // Get post by slug

	// Protected routes
	posts.Use(authMiddleware)
	posts.Post("/", postHandler.CreatePost)                  // Create new post
	posts.Get("/my-posts", postHandler.ListUserPosts)        // Get current user's posts
	posts.Put("/:id", postHandler.UpdatePost)                // Update post
	posts.Delete("/:id", postHandler.DeletePost)             // Delete post
	posts.Patch("/:id/publish", postHandler.PublishPost)     // Publish post
	posts.Patch("/:id/unpublish", postHandler.UnpublishPost) // Unpublish post
}
