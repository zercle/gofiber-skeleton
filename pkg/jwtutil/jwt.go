package jwtutil

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Expiration time.Duration
}

func NewJWT(privateKeyBase64, publicKeyBase64 string) *JWT {
	privBytes, _ := base64.StdEncoding.DecodeString(privateKeyBase64)
	pubBytes, _ := base64.StdEncoding.DecodeString(publicKeyBase64)
	return &JWT{
		PrivateKey: ed25519.PrivateKey(privBytes),
		PublicKey:  ed25519.PublicKey(pubBytes),
		Expiration: 24 * time.Hour,
	}
}

func (j *JWT) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.Expiration).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(j.PrivateKey)
}

func (j *JWT) ValidateToken(tokenStr string) bool {
	if tokenStr == "" {
		return false
	}
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.PublicKey, nil
	})
	return err == nil
}

func (j *JWT) ExtractUserID(tokenStr string) (uint, error) {
	if tokenStr == "" {
		return 0, fmt.Errorf("token is empty")
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.PublicKey, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token claims")
	}
	uidVal, ok := claims["user_id"]
	if !ok {
		return 0, fmt.Errorf("user_id not found in token claims")
	}
	switch v := uidVal.(type) {
	case float64:
		return uint(v), nil
	case string:
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(id), nil
	default:
		return 0, fmt.Errorf("invalid user_id type in token claims")
	}
}
