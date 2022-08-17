package datasources

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	helpers "github.com/zercle/gofiber-helpers"
)

var (
	// JWTVerifyKey JWT public key
	JWTVerifyKey     *crypto.PublicKey
	JWTSignKey       *crypto.PrivateKey
	JWTSigningMethod *jwt.SigningMethod
)

// ValidationJWT JWT validation func
var ValidationJWT jwt.Keyfunc = func(token *jwt.Token) (publicKey interface{}, err error) {
	if JWTVerifyKey == nil {
		err = fiber.NewError(http.StatusNotFound, "JWTVerifyKey not init yet")
	}
	// debug
	// log.Printf("source: %+v\nvalue: %+v", helpers.WhereAmI(), *JWTVerifyKey)
	return *JWTVerifyKey, err
}

func JTWLocalKey(privateKeyPath, publicKeyPath string) (privateKey crypto.PrivateKey, publicKey crypto.PublicKey, signingMethod jwt.SigningMethod, err error) {

	if len(privateKeyPath) == 0 {
		privateKeyPath = viper.GetString("jwt.private")
	}
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return
	}

	var publicKeyFile []byte
	if len(publicKeyPath) != 0 {
		// publicKeyPath = viper.GetString("jwt.public")
		publicKeyFile, err = os.ReadFile(publicKeyPath)
		if err != nil {
			err = fiber.NewError(http.StatusInternalServerError, err.Error())
			return
		}
	}

	// EdDSA
	if privateKey, err = jwt.ParseEdPrivateKeyFromPEM(privateKeyFile); err == nil {
		if publicKey, err = jwt.ParseEdPublicKeyFromPEM(publicKeyFile); err != nil {
			publicKey = privateKey.(ed25519.PrivateKey).Public()
			err = nil
		}
		signingMethod = jwt.SigningMethodEdDSA
		return
	}
	// debug
	// log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)

	// ECDSA
	if privateKey, err = jwt.ParseECPrivateKeyFromPEM(privateKeyFile); err == nil {
		if publicKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyFile); err != nil {
			publicKey = privateKey.(*ecdsa.PrivateKey).Public()
			err = nil
		}
		switch privateKey.(*ecdsa.PrivateKey).Curve.Params().BitSize {
		case 256:
			signingMethod = jwt.SigningMethodES256
		case 384:
			signingMethod = jwt.SigningMethodES384
		case 512:
			signingMethod = jwt.SigningMethodES512
		}
		return
	}
	// debug
	// log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)

	// RSA
	if privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile); err == nil {
		if publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyFile); err != nil {
			publicKey = privateKey.(*rsa.PrivateKey).Public()
			err = nil
		}
		switch privateKey.(*rsa.PrivateKey).N.BitLen() {
		case 256:
			signingMethod = jwt.SigningMethodRS256
		case 384:
			signingMethod = jwt.SigningMethodRS384
		case 512:
			signingMethod = jwt.SigningMethodRS512
		}
		return
	}
	// debug
	// log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)

	signingMethod = jwt.SigningMethodNone
	return
}
