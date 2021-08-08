package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/internal/datasources"
)

func ExtractBearerToken(authHeader string) (token string, err error) {
	authHeader = strings.TrimSpace(authHeader)
	authHeaders := strings.Split(authHeader, " ")
	if len(authHeaders) != 2 || authHeaders[0] != "Bearer" {
		err = fiber.NewError(http.StatusUnauthorized, "Authorization: Bearer token")
		return
	}
	token = authHeaders[1]
	return
}

func ExtractSocketToken(authHeader string) (token string, err error) {
	authHeaders := strings.Split(authHeader, ",")
	if len(authHeaders) < 2 || authHeaders[0] != "Bearer" {
		err = fiber.NewError(http.StatusUnauthorized, "Sec-WebSocket-Protocol: Bearer, access_token")
		return
	}
	token = strings.TrimSpace(authHeaders[1])
	return
}

// ReqLineAuthHandler check session
func ReqAuthHandler(c *fiber.Ctx) (err error) {
	tokenStr, err := ExtractBearerToken(c.Get(fiber.HeaderAuthorization))
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, err.Error())
	}

	jwtToken, err := jwt.Parse(tokenStr, datasources.ValidationJWT)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, err.Error())
	}
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		c.Locals("claims", claims)
	} else {
		// debug
		log.Printf("%+v\nvalue: %+v", helpers.WhereAmI(), claims)
		log.Printf("%+v\nvalue: %+v", helpers.WhereAmI(), jwtToken)
		return fiber.NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}
	return c.Next()
}
