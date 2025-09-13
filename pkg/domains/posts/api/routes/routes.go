package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zercle/gofiber-skeleton/pkg/domains/posts/api/handlers"
)

type PostRoutes struct {
	postHandler handlers.PostHandler
}

func NewPostRoutes(postHandler handlers.PostHandler) PostRoutes {
	return PostRoutes{postHandler: postHandler}
}

// RegisterRoutes registers all posts module routes.
func (r *PostRoutes) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	posts := router.Group("/posts")

	// Public routes
	posts.Get("/", r.postHandler.ListPosts)               // List all posts (can filter by published)
	posts.Get("/:id", r.postHandler.GetPost)              // Get post by ID
	posts.Get("/slug/:slug", r.postHandler.GetPostBySlug) // Get post by slug

	// Protected routes
	posts.Use(authMiddleware)
	posts.Post("/", r.postHandler.CreatePost)                  // Create new post
	posts.Get("/my-posts", r.postHandler.ListUserPosts)        // Get current user's posts
	posts.Put("/:id", r.postHandler.UpdatePost)                // Update post
	posts.Delete("/:id", r.postHandler.DeletePost)             // Delete post
	posts.Patch("/:id/publish", r.postHandler.PublishPost)     // Publish post
	posts.Patch("/:id/unpublish", r.postHandler.UnpublishPost) // Unpublish post
}