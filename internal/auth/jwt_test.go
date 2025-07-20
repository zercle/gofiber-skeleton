package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	userID := "test-user-id"
	secret := "test-secret"
	expiry := time.Hour

	tokenString, err := GenerateToken(userID, secret, expiry)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestValidateToken(t *testing.T) {
	userID := "test-user-id"
	secret := "test-secret"
	expiry := time.Hour

	tokenString, err := GenerateToken(userID, secret, expiry)
	assert.NoError(t, err)

	claims, err := ValidateToken(tokenString, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)

	// Test with expired token
	expiredTokenString, err := GenerateToken(userID, secret, -time.Hour)
	assert.NoError(t, err)

	_, err = ValidateToken(expiredTokenString, secret)
	assert.Error(t, err)

	// Test with invalid secret
	_, err = ValidateToken(tokenString, "wrong-secret")
	assert.Error(t, err)
}
