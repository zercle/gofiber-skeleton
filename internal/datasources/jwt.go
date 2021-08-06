package datasources

import (
	"crypto/ecdsa"
	"log"
	"net/http"
	"os"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

var (
	// JWTVerifyKey JWT public key
	JWTVerifyKey *ecdsa.PublicKey
	JWTSignKey   *ecdsa.PrivateKey
)

// ValidationJWT JWT validation func
var ValidationJWT jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
	return JWTVerifyKey, nil
}

func JTWLocalKey() (err error) {
	publicKey, err := os.ReadFile(os.Getenv("JWT_PUBLIC"))
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return err
	}
	JWTVerifyKey, err = jwt.ParseECPublicKeyFromPEM(publicKey)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return err
	}
	privateKey, err := os.ReadFile(os.Getenv("JWT_PRIVATE"))
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return err
	}
	JWTSignKey, err = jwt.ParseECPrivateKeyFromPEM(privateKey)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return err
	}
	return
}
