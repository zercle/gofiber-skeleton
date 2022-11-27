package datasources

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	helpers "github.com/zercle/gofiber-helpers"
)

type JwtResources struct {
	JwtVerifyKey     crypto.PublicKey
	JwtSignKey       crypto.PrivateKey
	JwtSigningMethod jwt.SigningMethod
}

func JTWLocalKey(privateKeyPath, publicKeyPath string) (jwtResources JwtResources, err error) {

	if len(privateKeyPath) == 0 {
		privateKeyPath = viper.GetString("jwt.private")
	}
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return
	}

	if len(publicKeyPath) == 0 {
		publicKeyPath = viper.GetString("jwt.public")
	}
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return
	}

	// EdDSA
	if jwtResources.JwtSignKey, err = jwt.ParseEdPrivateKeyFromPEM(privateKeyFile); err == nil {
		if jwtResources.JwtVerifyKey, err = jwt.ParseEdPublicKeyFromPEM(publicKeyFile); err != nil {
			jwtResources.JwtVerifyKey = jwtResources.JwtSignKey.(ed25519.PrivateKey).Public()
		}
		jwtResources.JwtSigningMethod = jwt.SigningMethodEdDSA
		return
	}

	// ECDSA
	if jwtResources.JwtSignKey, err = jwt.ParseECPrivateKeyFromPEM(privateKeyFile); err == nil {
		if jwtResources.JwtVerifyKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyFile); err != nil {
			log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
			jwtResources.JwtVerifyKey = jwtResources.JwtSignKey.(*ecdsa.PrivateKey).PublicKey
		}
		switch jwtResources.JwtSignKey.(*ecdsa.PrivateKey).Curve.Params().BitSize {
		case 256:
			jwtResources.JwtSigningMethod = jwt.SigningMethodES256
		case 384:
			jwtResources.JwtSigningMethod = jwt.SigningMethodES384
		case 512:
			jwtResources.JwtSigningMethod = jwt.SigningMethodES512
		}
		return
	}

	// RSA
	if jwtResources.JwtSignKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile); err == nil {
		if jwtResources.JwtVerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyFile); err != nil {
			log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
			jwtResources.JwtVerifyKey = jwtResources.JwtSignKey.(*rsa.PrivateKey).PublicKey
		}
		switch jwtResources.JwtSignKey.(*rsa.PrivateKey).N.BitLen() {
		case 256:
			jwtResources.JwtSigningMethod = jwt.SigningMethodRS256
		case 384:
			jwtResources.JwtSigningMethod = jwt.SigningMethodRS384
		case 512:
			jwtResources.JwtSigningMethod = jwt.SigningMethodRS512
		}
		return
	}

	if jwtResources.JwtSigningMethod == nil {
		jwtResources.JwtSigningMethod = jwt.SigningMethodNone
	}

	return
}

func InitJwtParser() (*jwt.Parser) {
	return jwt.NewParser()
}

func (r *JwtResources) SetKeyfunc() {
	JwtKeyfunc = func(token *jwt.Token) (publicKey interface{}, err error) {
		if r.JwtVerifyKey == nil {
			err = fiber.NewError(http.StatusFailedDependency, "JWTVerifyKey not init yet")
		}
		// debug
		// log.Printf("source: %+v\nvalue: %+v", helpers.WhereAmI(), *JWTVerifyKey)
		return r.JwtVerifyKey, err
	}
}

func (r *JwtResources) IsJwtActive(tokenStr string) (token *jwt.Token, isActive bool, err error) {
	if JwtParser == nil {
		err = fiber.NewError(http.StatusFailedDependency, "JwtParser not init yet")
		return
	}
	claims := jwt.RegisteredClaims{}
	token, _, err = JwtParser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		return
	}

	isActive = claims.VerifyExpiresAt(time.Now(), true)

	return
}
