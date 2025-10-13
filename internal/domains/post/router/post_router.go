package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/handler"
	"github.com/zercle/gofiber-skeleton/internal/shared/middleware"
	sharedRouter "github.com/zercle/gofiber-skeleton/internal/shared/router"
)

type PostRouter struct {
	*sharedRouter.BaseRouter
	postHandler  *handler.PostHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewPostRouter(i *do.Injector) (*PostRouter, error) {
	postHandler := do.MustInvoke[*handler.PostHandler](i)
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](i)

	return &PostRouter{
		BaseRouter:   sharedRouter.NewBaseRouter("/posts"),
		postHandler:  postHandler,
		authMiddleware: authMiddleware,
	}, nil
}

func (r *PostRouter) RegisterRoutes(app fiber.Router) {
	posts := app.Group(r.Prefix, r.authMiddleware.RequireAuth())
	posts.Post("/", r.postHandler.CreatePost)
	posts.Get("/", r.postHandler.ListPosts)
	posts.Get("/my", r.postHandler.ListMyPosts)
	posts.Get("/:id", r.postHandler.GetPost)
	posts.Put("/:id", r.postHandler.UpdatePost)
	posts.Delete("/:id", r.postHandler.DeletePost)
}