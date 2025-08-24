package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	sqlc "github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
	userhandler "github.com/zercle/gofiber-skeleton/internal/user/handler"
	userrepository "github.com/zercle/gofiber-skeleton/internal/user/repository"
	userusecase "github.com/zercle/gofiber-skeleton/internal/user/usecase"
)

func setupUserIntegrationTest(t *testing.T) (*fiber.App, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlcQueries := sqlc.New(db)
	userRepo := userrepository.NewUserRepository(sqlcQueries)
	userUseCase := userusecase.NewUserUseCase(userRepo, "test-jwt-secret") // Using a test secret
	userHandler := userhandler.NewUserHandler(userUseCase)

	app := fiber.New()
	app.Post("/api/v1/register", userHandler.Register)
	app.Post("/api/v1/login", userHandler.Login)
	// Add other user routes as needed, e.g., GetUser, UpdateUser, DeleteUser

	return app, mock, db
}

func TestUserIntegration_Register(t *testing.T) {
	app, mock, db := setupUserIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	registerInput := userhandler.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	expectedUserUUID := uuid.New()
	mockTime := time.Now()

	t.Run("successful user registration", func(t *testing.T) {

		// Mock query to check if username exists
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = $1`,
		)).
			WithArgs(registerInput.Username).
			WillReturnError(sql.ErrNoRows) // Simulate user not found

		// Mock insert user
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(expectedUserUUID, registerInput.Username, "hashed_password", domain.RoleCustomer, mockTime, mockTime)
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id, username, password_hash, role, created_at, updated_at`,
		)).
			WithArgs(registerInput.Username, sqlmock.AnyArg(), domain.RoleCustomer).
			WillReturnRows(rows)

		body, _ := json.Marshal(registerInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		data := responseBody["data"].(map[string]any)
		assert.Equal(t, "User registered successfully", data["message"])
		user := data["user"].(map[string]any)
		assert.Equal(t, registerInput.Username, user["username"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("registration with existing username", func(t *testing.T) {

		// Mock query to check if username exists and return existing user
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(uuid.New(), registerInput.Username, "hashedpassword", domain.RoleCustomer, mockTime, mockTime)
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = $1`,
		)).
			WithArgs(registerInput.Username).
			WillReturnRows(rows) // Simulate user found

		body, _ := json.Marshal(registerInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "fail", responseBody["status"])
		assert.Equal(t, "username already exists", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserIntegration_Login(t *testing.T) {
	app, mock, db := setupUserIntegrationTest(t)
	defer func() {
		_ = db.Close()
	}()

	loginInput := userhandler.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	expectedUserUUID := uuid.New()
	mockTime := time.Now()

	t.Run("successful user login", func(t *testing.T) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginInput.Password), bcrypt.DefaultCost)
		require.NoError(t, err)

		// Mock query to find user by username
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(expectedUserUUID, loginInput.Username, string(hashedPassword), domain.RoleCustomer, mockTime, mockTime) // Dynamically generated bcrypt hash
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = $1`,
		)).
			WithArgs(loginInput.Username).
			WillReturnRows(rows)

		body, _ := json.Marshal(loginInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "success", responseBody["status"])
		assert.Contains(t, responseBody, "data")
		data := responseBody["data"].(map[string]any)
		assert.Equal(t, "Login successful", data["message"])
		assert.Contains(t, data, "token")
		assert.Contains(t, data, "user")
		user := data["user"].(map[string]any)
		assert.Equal(t, loginInput.Username, user["username"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("login with invalid credentials", func(t *testing.T) {

		// Mock query to find user by username (simulating user not found or wrong password check)
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = $1`,
		)).
			WithArgs(loginInput.Username).
			WillReturnError(sql.ErrNoRows) // Simulate user not found

		body, _ := json.Marshal(loginInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "invalid credentials", responseBody["message"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
