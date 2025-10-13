package router

import "github.com/gofiber/fiber/v2"

// Router defines the interface for domain-specific routers.
type Router interface {
	RegisterRoutes(app fiber.Router)
}

// BaseRouter provides common functionality for domain routers.
type BaseRouter struct {
	Prefix string
}

// NewBaseRouter creates a new BaseRouter instance.
func NewBaseRouter(prefix string) *BaseRouter {
	return &BaseRouter{
		Prefix: prefix,
	}
}

// Group creates a new Fiber group with the router's prefix.
func (r *BaseRouter) Group(app fiber.Router, path string, handlers ...fiber.Handler) fiber.Router {
	return app.Group(r.Prefix+path, handlers...)
}