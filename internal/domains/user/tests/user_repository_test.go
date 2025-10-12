package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/mocks"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	"github.com/zercle/gofiber-skeleton/pkg/database"
)

func TestUserRepository_Create_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Mock database expectations
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs("test@example.com", "hashed_password", "Test User").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Test data
	user := entity.NewUser("test@example.com", "hashed_password", "Test User")

	// Execute
	err = repo.Create(context.Background(), user)

	// Assert
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "is_active", "created_at", "updated_at"}).
		AddRow(userID.String(), "test@example.com", "hashed_password", "Test User", true, time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, is_active, created_at, updated_at FROM users`).
		WithArgs(userID.String()).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetByID(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.FullName)
	assert.True(t, result.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, is_active, created_at, updated_at FROM users`).
		WithArgs(userID.String()).
		WillReturnError(sql.ErrNoRows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetByID(context.Background(), userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "is_active", "created_at", "updated_at"}).
		AddRow(uuid.New().String(), "test@example.com", "hashed_password", "Test User", true, time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, is_active, created_at, updated_at FROM users`).
		WithArgs("test@example.com").
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetByEmail(context.Background(), "test@example.com")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.FullName)
	assert.True(t, result.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, email, password_hash, full_name, is_active, created_at, updated_at FROM users`).
		WithArgs("nonexistent@example.com").
		WillReturnError(sql.ErrNoRows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetByEmail(context.Background(), "nonexistent@example.com")

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`UPDATE users`).
		WithArgs("updated@example.com", "Updated User", gomock.Any()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Test data
	user := &entity.DomainUser{
		ID:        userID,
		Email:     "updated@example.com",
		FullName:  "Updated User",
		UpdatedAt: time.Now(),
	}

	// Execute
	err = repo.Update(context.Background(), user)

	// Assert
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`UPDATE users`).
		WithArgs("test@example.com", "Test User", gomock.Any()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Test data
	user := &entity.DomainUser{
		ID:        userID,
		Email:     "test@example.com",
		FullName:  "Test User",
		UpdatedAt: time.Now(),
	}

	// Execute
	err = repo.Update(context.Background(), user)

	// Assert
	require.Error(t, err)
	assert.Equal(t, entity.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdatePassword_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`UPDATE users`).
		WithArgs(gomock.Any(), "new_hashed_password").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	err = repo.UpdatePassword(context.Background(), userID, "new_hashed_password")

	// Assert
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdatePassword_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`UPDATE users`).
		WithArgs(gomock.Any(), "new_hashed_password").
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	err = repo.UpdatePassword(context.Background(), userID, "new_hashed_password")

	// Assert
	require.Error(t, err)
	assert.Equal(t, entity.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Deactivate_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`UPDATE users`).
		WithArgs(userID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	err = repo.Deactivate(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Deactivate_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`UPDATE users`).
		WithArgs(userID.String()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	err = repo.Deactivate(context.Background(), userID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, entity.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_List_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "full_name", "is_active", "created_at", "updated_at"}).
		AddRow(uuid.New().String(), "user1@example.com", "User One", true, time.Now(), time.Now()).
		AddRow(uuid.New().String(), "user2@example.com", "User Two", true, time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, email, full_name, is_active, created_at, updated_at FROM users`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.List(context.Background(), 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "user1@example.com", result[0].Email)
	assert.Equal(t, "user2@example.com", result[1].Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_EmailExists_True(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Mock database expectations
	mock.ExpectQuery(`SELECT COUNT`).
		WithArgs("existing@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	exists, err := repo.EmailExists(context.Background(), "existing@example.com")

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_EmailExists_False(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Mock database expectations
	mock.ExpectQuery(`SELECT COUNT`).
		WithArgs("nonexistent@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	exists, err := repo.EmailExists(context.Background(), "nonexistent@example.com")

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Count_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Mock database expectations
	mock.ExpectQuery(`SELECT COUNT`).
		WithArgs("").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewUserRepository(injector)
	require.NoError(t, err)

	// Execute
	count, err := repo.Count(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 5, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
