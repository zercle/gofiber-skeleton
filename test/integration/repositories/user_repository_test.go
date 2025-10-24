package repositories_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zercle/template-go-fiber/internal/domains"
	"github.com/zercle/template-go-fiber/internal/repositories"
	db "github.com/zercle/template-go-fiber/internal/infrastructure/sqlc"
)

func TestUserRepository_GetByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	// Mock the database query
	rows := sqlmock.NewRows([]string{
		"id", "email", "password_hash", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		"user-123",
		"test@example.com",
		"hashed_password",
		"John",
		"Doe",
		true,
		time.Now(),
		time.Now(),
		nil,
	)

	mock.ExpectQuery("SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at FROM users WHERE id = .* AND deleted_at IS NULL").
		WithArgs("user-123").
		WillReturnRows(rows)

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	user, err := repo.GetByID(context.Background(), "user-123")

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.ID != "user-123" {
		t.Errorf("expected ID user-123, got %s", user.ID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	// Mock no rows returned
	mock.ExpectQuery("SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at FROM users WHERE id = .* AND deleted_at IS NULL").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	user, err := repo.GetByID(context.Background(), "nonexistent")

	// sqlmock returns error, but our code should handle it
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("unexpected error: %v", err)
	}

	if user != nil {
		t.Errorf("expected nil user for nonexistent ID, got %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	_ = err // Handle err

	// Mock the INSERT query - Note: created_at and updated_at use NOW() in SQL, not parameters
	mock.ExpectExec("INSERT INTO users \\(id, email, password_hash, first_name, last_name, is_active, created_at, updated_at\\)").
		WithArgs("new-user", "new@example.com", "hashed", "Jane", "Smith", true).
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	user := &domains.User{
		ID:           "new-user",
		Email:        "new@example.com",
		PasswordHash: "hashed",
		FirstName:    stringPtr("Jane"),
		LastName:     stringPtr("Smith"),
		IsActive:     true,
	}

	err = repo.Create(context.Background(), user)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	// Mock the UPDATE query
	mock.ExpectExec("UPDATE users SET email = .*, password_hash = .*, first_name = .*, last_name = .*, is_active = .*, updated_at = NOW\\(\\) WHERE id = .* AND deleted_at IS NULL").
		WithArgs("updated@example.com", "new_hash", "Jane", "Smith", true, "user-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	user := &domains.User{
		ID:           "user-123",
		Email:        "updated@example.com",
		PasswordHash: "new_hash",
		FirstName:    stringPtr("Jane"),
		LastName:     stringPtr("Smith"),
		IsActive:     true,
	}

	err = repo.Update(context.Background(), user)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	// Mock the soft DELETE query
	mock.ExpectExec("UPDATE users SET deleted_at = NOW\\(\\), updated_at = NOW\\(\\) WHERE id = .* AND deleted_at IS NULL").
		WithArgs("user-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	err = repo.Delete(context.Background(), "user-123")

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_List(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	// Mock the SELECT query for list
	rows := sqlmock.NewRows([]string{
		"id", "email", "password_hash", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at",
	}).
		AddRow(
			"user-1",
			"user1@example.com",
			"hash1",
			"John",
			"Doe",
			true,
			time.Now(),
			time.Now(),
			nil,
		).
		AddRow(
			"user-2",
			"user2@example.com",
			"hash2",
			"Jane",
			"Smith",
			true,
			time.Now(),
			time.Now(),
			nil,
		)

	mock.ExpectQuery("SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT .* OFFSET .*").
		WithArgs(int64(10), int64(0)).
		WillReturnRows(rows)

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	users, err := repo.List(context.Background(), 10, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_Count(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	// Mock the COUNT query
	rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(5))

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) as count FROM users WHERE deleted_at IS NULL").
		WillReturnRows(rows)

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	count, err := repo.Count(context.Background())

	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if count != 5 {
		t.Errorf("expected count 5, got %d", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer func() { _ = mockDB.Close() }()

	// Mock the SELECT by email query
	rows := sqlmock.NewRows([]string{
		"id", "email", "password_hash", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		"user-123",
		"test@example.com",
		"hashed_password",
		"John",
		"Doe",
		true,
		time.Now(),
		time.Now(),
		nil,
	)

	mock.ExpectQuery("SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at FROM users WHERE email = .* AND deleted_at IS NULL").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	queries := db.New(mockDB)
	repo := repositories.NewUserRepository(queries)

	user, err := repo.GetByEmail(context.Background(), "test@example.com")

	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
