package routers

import (
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

func (r *RouterResources) Index() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
}
