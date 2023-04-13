package datasources

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"
)

func NewJwtLocalKey(privateKeyPath string) (jwtSignKey crypto.PrivateKey, jwtVerifyKey crypto.PublicKey, jwtSigningMethod jwt.SigningMethod, err error) {

	jwtSigningMethod = jwt.SigningMethodNone

	if len(privateKeyPath) == 0 {
		err = fmt.Errorf("InitJwtLocalKey: %+s", "need privateKeyPath")
		return
	}
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("InitJwtLocalKey: %+v", err)
		return
	}

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(privateKeyFile); block == nil {
		err = jwt.ErrKeyMustBePEMEncoded
		log.Printf("InitJwtLocalKey: %+v", err)
		return
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return
		}
	}

	// determind which signing method
	switch key := parsedKey.(type) {
	case ed25519.PrivateKey:
		jwtSignKey = key
		jwtVerifyKey = key.Public()
		jwtSigningMethod = jwt.SigningMethodEdDSA
	case *ecdsa.PrivateKey:
		jwtSignKey = key
		jwtVerifyKey = key.Public()
		switch key.Curve {
		case elliptic.P256():
			jwtSigningMethod = jwt.SigningMethodES256
		case elliptic.P384():
			jwtSigningMethod = jwt.SigningMethodES384
		case elliptic.P521():
			jwtSigningMethod = jwt.SigningMethodES512
		}
	case *rsa.PrivateKey:
		jwtSignKey = key
		jwtVerifyKey = key.Public()
		switch key.Size() {
		case 256:
			jwtSigningMethod = jwt.SigningMethodRS256
		case 384:
			jwtSigningMethod = jwt.SigningMethodRS384
		case 512:
			jwtSigningMethod = jwt.SigningMethodRS512
		}
	default:
		err = errors.New("unsupported private key type")
		return
	}

	return
}
