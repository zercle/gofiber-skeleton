package datasources

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	helpers "github.com/zercle/gofiber-helpers"
)

type JwtResources struct {
	JwtVerifyKey     crypto.PublicKey
	JwtSignKey       crypto.PrivateKey
	JwtSigningMethod jwt.SigningMethod
	JwtKeyfunc       jwt.Keyfunc
	JwtParser        *jwt.Parser
}

func NewJWT(privateKeyPath, publicKeyPath string) (jwtResources *JwtResources, err error) {
	resources, err := JTWLocalKey(privateKeyPath, publicKeyPath)
	resources.JwtKeyfunc = func(token *jwt.Token) (publicKey interface{}, err error) {
		if resources.JwtVerifyKey == nil {
			err = fmt.Errorf("JWTVerifyKey not init yet")
		}
		return resources.JwtVerifyKey, err
	}
	resources.JwtParser = jwt.NewParser()
	return &resources, err
}

func JTWLocalKey(privateKeyPath, publicKeyPath string) (jwtResources JwtResources, err error) {

	if len(privateKeyPath) == 0 {
		privateKeyPath = viper.GetString("jwt.private")
	}
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		return
	}

	if len(publicKeyPath) == 0 {
		publicKeyPath = viper.GetString("jwt.public")
	}
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
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

func (r *JwtResources) IsJwtActive(tokenStr string) (token *jwt.Token, isActive bool, err error) {
	if r.JwtParser == nil {
		err = fmt.Errorf("JwtParser not init yet")
		return
	}
	claims := jwt.RegisteredClaims{}
	token, _, err = r.JwtParser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		log.Printf("source: %+v\nerr: %+v", helpers.WhereAmI(), err)
		return
	}

	isActive = claims.VerifyExpiresAt(time.Now(), true)

	return
}
