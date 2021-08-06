package routes

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	helpers "github.com/zercle/gofiber-helpers"
)

// Index route
func Index(c *fiber.Ctx) (err error) {
	responseForm := helpers.ResponseForm{
		Success: true,
		Result: fiber.Map{
			"ip":            c.IP(),
			"remote":        c.IPs(),
			"client":        string(c.Context().UserAgent()),
			"secure":        c.Secure(),
			"conn_time":     c.Context().ConnTime(),
			"response_time": c.Context().Time(),
			"version":       c.Params("version", "0"),
		},
	}
	return c.JSON(responseForm)
}

func (r *RouterResources) SetRouters(app *fiber.App) {
	apiLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get(fiber.HeaderXForwardedFor)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
		},
	})

	app.Get("/", Index)

	apiv1Group := app.Group("/api/v1", apiLimiter)
	{
		apiv1Group.Get("/", Index)
	}
}
