package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/db"
	"github.com/zercle/gofiber-skeleton/internal/user/repository"
)

func TestUserRepository_CreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewUserRepository(queries)

	ctx := context.Background()
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashedpassword"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(userID, username, email, passwordHash, now, now)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(username, email, passwordHash).
		WillReturnRows(rows)

	user, err := repo.Create(ctx, username, email, passwordHash)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewUserRepository(queries)

	ctx := context.Background()
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashedpassword"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(userID, username, email, passwordHash, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(ctx, email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewUserRepository(queries)

	ctx := context.Background()
	email := "nonexistent@example.com"

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByEmail(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewUserRepository(queries)

	ctx := context.Background()
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashedpassword"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(userID, username, email, passwordHash, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE id`).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, username, user.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	queries := db.New(mockDB)
	repo := repository.NewUserRepository(queries)

	ctx := context.Background()
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashedpassword"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(userID, username, email, passwordHash, now, now)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE username`).
		WithArgs(username).
		WillReturnRows(rows)

	user, err := repo.GetByUsername(ctx, username)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
