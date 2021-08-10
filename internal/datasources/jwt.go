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
	JWTVerifyKey     *ecdsa.PublicKey
	JWTSignKey       *ecdsa.PrivateKey
	JWTSigningMethod *jwt.SigningMethodECDSA
)

// ValidationJWT JWT validation func
var ValidationJWT jwt.Keyfunc = func(token *jwt.Token) (publicKey interface{}, err error) {
	if JWTVerifyKey == nil {
		err = fiber.NewError(http.StatusNotFound, "JWTVerifyKey not init yet")
	}
	return JWTVerifyKey, err
}

func JTWLocalKey(privateKeyPath, publicKeyPath string) (privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, signingMethod *jwt.SigningMethodECDSA, err error) {

	if len(privateKeyPath) == 0 {
		privateKeyPath = os.Getenv("JWT_PRIVATE")
	}
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return
	}
	privateKey, err = jwt.ParseECPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return
	}

	switch privateKey.Curve.Params().BitSize {
	case 256:
		signingMethod = jwt.SigningMethodES256
	case 384:
		signingMethod = jwt.SigningMethodES384
	case 512:
		signingMethod = jwt.SigningMethodES512
	}

	if len(publicKeyPath) == 0 {
		publicKeyPath = os.Getenv("JWT_PUBLIC")
	}
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		publicKey = &privateKey.PublicKey
		// log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		// err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return
	}
	publicKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return
	}

	return
}
