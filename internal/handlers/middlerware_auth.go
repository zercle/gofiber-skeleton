package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v4"
	helpers "github.com/zercle/gofiber-helpers"
)

var ApiLimiter = limiter.New(limiter.Config{
	Max:        750,
	Expiration: 30 * time.Second,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.Get(fiber.HeaderXForwardedFor)
	},
	LimitReached: func(c *fiber.Ctx) error {
		return helpers.NewError(http.StatusTooManyRequests, helpers.WhereAmI(), http.StatusText(http.StatusTooManyRequests))
	},
})

func ExtractBearerToken(authHeader string) (token string, err error) {
	authHeader = strings.TrimSpace(authHeader)
	authHeaders := strings.Split(authHeader, " ")
	if len(authHeaders) != 2 || authHeaders[0] != "Bearer" {
		err = helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), "Authorization: Bearer token")
		return
	}
	token = authHeaders[1]
	return
}

func ExtractSocketToken(authHeader string) (token string, err error) {
	authHeaders := strings.Split(authHeader, ",")
	if len(authHeaders) < 2 || authHeaders[0] != "Bearer" {
		err = helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), "Sec-WebSocket-Protocol: Bearer, access_token")
		return
	}
	token = strings.TrimSpace(authHeaders[1])
	return
}

func ExtractLevel(aud []string) (level int, err error) {
	if len(aud) < 1 {
		err = helpers.NewError(http.StatusBadRequest, helpers.WhereAmI(), "aud field missmatch")
		return
	}
	levels := strings.Split(aud[0], ":")
	if len(levels) < 1 {
		err = helpers.NewError(http.StatusBadRequest, helpers.WhereAmI(), "levels field missmatch")
		return
	}

	return strconv.Atoi(levels[1])
}

// ReqLineAuthHandler check session
func (r *RouterResources) ReqAuthHandler(reqLevels ...int) fiber.Handler {
	reqLevel := 4
	if len(reqLevels) != 0 {
		reqLevel = reqLevels[0]
	}

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		tokenStr, err := ExtractBearerToken(c.Get(fiber.HeaderAuthorization))
		if err != nil {
			return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), err.Error())
		}

		claims := new(jwt.RegisteredClaims)
		jwtToken, err := jwt.ParseWithClaims(tokenStr, claims, r.JwtResources.JwtKeyfunc)
		if err != nil {
			return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), err.Error())
		}
		if jwtToken != nil && jwtToken.Valid {
			if level, err := ExtractLevel(claims.Audience); err != nil {
				return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), err.Error())
			} else if level < reqLevel {
				return helpers.NewError(http.StatusForbidden, helpers.WhereAmI(), fmt.Sprintf("%s need permission level %d", c.Route().Path, reqLevel))
			} else {
				c.Locals("level", level)
			}
			c.Locals("claims", claims)
			c.Locals("token", jwtToken)
		} else {
			// debug
			log.Printf("%+v\nvalue: %+v", helpers.WhereAmI(), claims)
			log.Printf("%+v\nvalue: %+v", helpers.WhereAmI(), jwtToken)
			return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), http.StatusText(http.StatusUnauthorized))
		}
		return c.Next()
	}
}
