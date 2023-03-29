package infrastructure

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"log"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"
)

func NewJwtLocalKey(privateKeyPath, publicKeyPath string) (jwtSignKey crypto.PrivateKey, jwtVerifyKey crypto.PrivateKey, jwtSigningMethod jwt.SigningMethod, err error) {

	if len(privateKeyPath) == 0 {
		err = fmt.Errorf("InitJwtLocalKey: %+s", "need privateKeyPath")
		return
	}
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("InitJwtLocalKey: %+v", err)
		return
	}

	var publicKeyFile []byte
	if len(publicKeyPath) == 0 {
		log.Printf("InitJwtLocalKey: %+v", "empty publicKeyPath")
	} else {
		publicKeyFile, err = os.ReadFile(publicKeyPath)
		if err != nil {
			log.Printf("InitJwtLocalKey: %+v", err)
			return
		}
	}

	// EdDSA
	if jwtSignKey, err = jwt.ParseEdPrivateKeyFromPEM(privateKeyFile); err == nil {
		if jwtVerifyKey, err = jwt.ParseEdPublicKeyFromPEM(publicKeyFile); err != nil {
			jwtVerifyKey = jwtSignKey.(ed25519.PrivateKey).Public()
		}
		jwtSigningMethod = jwt.SigningMethodEdDSA
		return
	}

	// ECDSA
	if jwtSignKey, err = jwt.ParseECPrivateKeyFromPEM(privateKeyFile); err == nil {
		if jwtVerifyKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyFile); err != nil {
			log.Printf("InitJwtLocalKey: %+v", err)
			jwtVerifyKey = jwtSignKey.(*ecdsa.PrivateKey).PublicKey
		}
		switch jwtSignKey.(*ecdsa.PrivateKey).Curve.Params().BitSize {
		case 256:
			jwtSigningMethod = jwt.SigningMethodES256
		case 384:
			jwtSigningMethod = jwt.SigningMethodES384
		case 512:
			jwtSigningMethod = jwt.SigningMethodES512
		}
		return
	}

	// RSA
	if jwtSignKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile); err == nil {
		if jwtVerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyFile); err != nil {
			log.Printf("InitJwtLocalKey: %+v", err)
			jwtVerifyKey = jwtSignKey.(*rsa.PrivateKey).PublicKey
		}
		switch jwtSignKey.(*rsa.PrivateKey).N.BitLen() {
		case 256:
			jwtSigningMethod = jwt.SigningMethodRS256
		case 384:
			jwtSigningMethod = jwt.SigningMethodRS384
		case 512:
			jwtSigningMethod = jwt.SigningMethodRS512
		}
		return
	}

	if jwtSigningMethod == nil {
		jwtSigningMethod = jwt.SigningMethodNone
	}

	return
}
