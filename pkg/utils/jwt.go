package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
)

type JWTManager struct {
	config *config.Config
}

func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{config: cfg}
}

func (j *JWTManager) GenerateToken(userID, email string) (string, error) {
	now := time.Now()
	
	claims := &middleware.AuthClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    j.config.JWT.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{j.config.JWT.Issuer},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.JWT.ExpiresIn)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.JWT.Secret))
}

func (j *JWTManager) ValidateToken(tokenString string) (*middleware.AuthClaims, error) {
	claims := &middleware.AuthClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}