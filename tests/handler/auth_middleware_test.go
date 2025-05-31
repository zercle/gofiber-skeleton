package handler_test

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gofiber-skeleton/internal/middleware"
	"gofiber-skeleton/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

func generateEd25519Keys(t *testing.T) (privStr, pubStr string, priv ed25519.PrivateKey, pub ed25519.PublicKey) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("failed to generate ed25519 keys: %v", err)
	}
	privStr = base64.StdEncoding.EncodeToString(priv)
	pubStr = base64.StdEncoding.EncodeToString(pub)
	return privStr, pubStr, priv, pub
}

func makeJWT(t *testing.T, jwt *jwtutil.JWT, userID uint) string {
	token, err := jwt.GenerateToken(userID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

func setupAuthApp(jwt *jwtutil.JWT) *fiber.App {
	app := fiber.New()
	app.Use(middleware.AuthMiddleware(jwt))
	app.Get("/protected", func(c *fiber.Ctx) error {
		uid := c.Locals(middleware.UserIDContextKey)
		return c.JSON(fiber.Map{"user_id": uid})
	})
	return app
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	privStr, pubStr, _, _ := generateEd25519Keys(t)
	jwt := jwtutil.NewJWT(privStr, pubStr)
	token := makeJWT(t, jwt, 42)

	app := setupAuthApp(jwt)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	buf := make([]byte, 128)
	n, _ := resp.Body.Read(buf)
	body := string(buf[:n])
	if !strings.Contains(body, "42") {
		t.Errorf("expected user_id 42 in response, got %s", body)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	privStr, pubStr, _, _ := generateEd25519Keys(t)
	jwt := jwtutil.NewJWT(privStr, pubStr)
	app := setupAuthApp(jwt)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	privStr, pubStr, _, _ := generateEd25519Keys(t)
	jwt := jwtutil.NewJWT(privStr, pubStr)
	app := setupAuthApp(jwt)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token something")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	privStr, pubStr, _, _ := generateEd25519Keys(t)
	jwt := jwtutil.NewJWT(privStr, pubStr)
	app := setupAuthApp(jwt)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Custom JWT with short expiration
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("failed to generate keys: %v", err)
	}
	privStr := base64.StdEncoding.EncodeToString(priv)
	pubStr := base64.StdEncoding.EncodeToString(pub)
	jwt := &jwtutil.JWT{
		PrivateKey: priv,
		PublicKey:  pub,
		Expiration: -1 * time.Second, // already expired
	}
	token, err := jwt.GenerateToken(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	app := setupAuthApp(jwtutil.NewJWT(privStr, pubStr))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidUserIDInToken(t *testing.T) {
	// Generate a token with user_id = 0
	privStr, pubStr, _, _ := generateEd25519Keys(t)
	jwt := jwtutil.NewJWT(privStr, pubStr)
	token := makeJWT(t, jwt, 0)
	app := setupAuthApp(jwt)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
