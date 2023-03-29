package infrastructure

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

var customErrorHandler = func(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := http.StatusInternalServerError
	var source string

	if e, ok := err.(*helpers.Error); ok {
		// Override status code if helpers.Error type
		code = e.Code
		source, _ = e.Source.(string)
	}

	responseData := helpers.ResponseForm{
		Success: false,
		Errors:  []helpers.ResponseError{helpers.ResponseError(*helpers.NewError(code, source, err.Error()))},
	}

	// Return statuscode with error message
	err = c.Status(code).JSON(responseData)
	if err != nil {
		// In case the JSON fails
		log.Printf("customErrorHandler: %+v", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Return from handler
	return nil
}
