package routes

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/zercle/gofiber-skelton/internal/datasources"
)

var apiLimiter = limiter.New(limiter.Config{
	Max:        750,
	Expiration: 30 * time.Second,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.Get(fiber.HeaderXForwardedFor)
	},
	LimitReached: func(c *fiber.Ctx) error {
		return fiber.NewError(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
	},
})

// RouterResources DB handler
type RouterResources struct {
	*datasources.Resources
}

// InitRouterResources returns a new DBHandler
func InitRouterResources(resources *datasources.Resources) *RouterResources {
	return &RouterResources{resources}
}
