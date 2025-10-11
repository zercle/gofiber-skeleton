package auth_test

import (
	"testing"

	"github.com/zercle/gofiber-skeleton/pkg/auth"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "TestPassword123"
	hash, _ := auth.HashPassword(password)

	t.Run("Correct Password", func(t *testing.T) {
		err := auth.CheckPassword(password, hash)
		if err != nil {
			t.Errorf("Expected no error for correct password, got %v", err)
		}
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		err := auth.CheckPassword("WrongPassword", hash)
		if err == nil {
			t.Error("Expected error for incorrect password")
		}
	})
}

func TestHashPasswordDeterminism(t *testing.T) {
	password := "TestPassword123"

	hash1, _ := auth.HashPassword(password)
	hash2, _ := auth.HashPassword(password)

	if hash1 == hash2 {
		t.Error("Expected different hashes for same password (bcrypt should use different salts)")
	}

	// Both hashes should validate the same password
	if err := auth.CheckPassword(password, hash1); err != nil {
		t.Error("First hash should validate password")
	}
	if err := auth.CheckPassword(password, hash2); err != nil {
		t.Error("Second hash should validate password")
	}
}
