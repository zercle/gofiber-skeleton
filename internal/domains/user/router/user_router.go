package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/handler"
	"github.com/zercle/gofiber-skeleton/internal/shared/middleware"
	sharedRouter "github.com/zercle/gofiber-skeleton/internal/shared/router"
)

type UserRouter struct {
	*sharedRouter.BaseRouter
	userHandler  *handler.UserHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewUserRouter(i *do.Injector) (*UserRouter, error) {
	userHandler := do.MustInvoke[*handler.UserHandler](i)
	authMiddleware := do.MustInvoke[*middleware.AuthMiddleware](i)

	return &UserRouter{
		BaseRouter:   sharedRouter.NewBaseRouter("/auth"),
		userHandler:  userHandler,
		authMiddleware: authMiddleware,
	}, nil
}

func (r *UserRouter) RegisterRoutes(app fiber.Router) {
	// Public routes
	auth := app.Group(r.Prefix)
	auth.Post("/register", r.userHandler.Register)
	auth.Post("/login", r.userHandler.Login)

	// Protected routes
	users := app.Group("/users", r.authMiddleware.RequireAuth())
	users.Get("/profile", r.userHandler.GetProfile)
	users.Put("/profile", r.userHandler.UpdateProfile)
	users.Delete("/profile", r.userHandler.DeleteAccount)
	users.Get("/:id", r.userHandler.GetByID)
	users.Get("/", r.userHandler.List)
}