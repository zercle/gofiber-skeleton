package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var _ JWTService = (*jwtService)(nil)

type JWTService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*jwtCustomClaims, error)
}

type jwtService struct {
	secretKey string
	issuer    string
	ttl       time.Duration
}

func NewJWTService(secretKey string, ttl time.Duration) JWTService {
	if ttl == 0 {
		ttl = 24 * time.Hour
	}

	return &jwtService{
		secretKey: secretKey,
		issuer:    "server",
		ttl:       ttl,
	}
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *jwtService) GenerateToken(userID string) (string, error) {
	currentTime := time.Now()
	claims := &jwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (*jwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
