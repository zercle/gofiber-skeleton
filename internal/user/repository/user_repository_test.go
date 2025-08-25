package userrepository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

func TestUserRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		_ = db.Close()
	}()

	mockQuerier := sqlc.New(db)
	repo := NewUserRepository(mockQuerier)

	user := &domain.User{
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "customer",
	}

	dbUser := sqlc.User{
		ID:           uuid.New(),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	t.Run("successful user creation", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users \\(username, password_hash, role\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id, username, password_hash, role, created_at, updated_at").
			WithArgs(user.Username, user.PasswordHash, user.Role).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
				AddRow(dbUser.ID, dbUser.Username, dbUser.PasswordHash, dbUser.Role, dbUser.CreatedAt.Time, dbUser.UpdatedAt.Time))

		createdUser, err := repo.Create(user)
		require.NoError(t, err)
		assert.Equal(t, dbUser.ID.String(), createdUser.ID)
		assert.False(t, createdUser.CreatedAt.IsZero())
		assert.False(t, createdUser.UpdatedAt.IsZero())
	})

	t.Run("database error on create", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users \\(username, password_hash, role\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id, username, password_hash, role, created_at, updated_at").
			WithArgs(user.Username, user.PasswordHash, user.Role).
			WillReturnError(errors.New("db error"))

		createdUser, err := repo.Create(user)
		assert.Nil(t, createdUser)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		_ = db.Close()
	}()

	mockQuerier := sqlc.New(db)
	repo := NewUserRepository(mockQuerier)

	userID := uuid.New().String()
	dbUser := sqlc.User{
		ID:           uuid.MustParse(userID),
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "customer",
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	t.Run("successful user retrieval by ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(dbUser.ID, dbUser.Username, dbUser.PasswordHash, dbUser.Role, dbUser.CreatedAt.Time, dbUser.UpdatedAt.Time)

		mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE id = \\$1").
			WithArgs(uuid.MustParse(userID)).
			WillReturnRows(rows)

		user, err := repo.GetByID(userID)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, dbUser.Username, user.Username)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		user, err := repo.GetByID("invalid-uuid")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid UUID length")
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE id = \\$1").
			WithArgs(uuid.MustParse(userID)).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("database error on get by ID", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE id = \\$1").
			WithArgs(uuid.MustParse(userID)).
			WillReturnError(errors.New("db error"))

		user, err := repo.GetByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		_ = db.Close()
	}()

	mockQuerier := sqlc.New(db)
	repo := NewUserRepository(mockQuerier)

	userID := uuid.New().String()
	dbUser := sqlc.User{
		ID:           uuid.MustParse(userID),
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Role:         "customer",
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	t.Run("successful user retrieval by username", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(dbUser.ID, dbUser.Username, dbUser.PasswordHash, dbUser.Role, dbUser.CreatedAt.Time, dbUser.UpdatedAt.Time)

		mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = \\$1").
			WithArgs("testuser").
			WillReturnRows(rows)

		user, err := repo.GetByUsername("testuser")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, userID, user.ID)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = \\$1").
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByUsername("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("database error on get by username", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = \\$1").
			WithArgs("testuser").
			WillReturnError(errors.New("db error"))

		user, err := repo.GetByUsername("testuser")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockQuerier := sqlc.New(db)
	repo := NewUserRepository(mockQuerier)

	userID := uuid.New().String()
	userToUpdate := &domain.User{
		ID:           userID,
		Username:     "updateduser",
		PasswordHash: "newhashedpassword",
		Role:         "admin",
	}

	dbUser := sqlc.User{
		ID:           uuid.MustParse(userID),
		Username:     userToUpdate.Username,
		PasswordHash: userToUpdate.PasswordHash,
		Role:         userToUpdate.Role,
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	t.Run("successful user update", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(dbUser.ID, dbUser.Username, dbUser.PasswordHash, dbUser.Role, dbUser.CreatedAt.Time, dbUser.UpdatedAt.Time)

		mock.ExpectQuery("UPDATE users SET username = \\$2, password_hash = \\$3, role = \\$4, updated_at = NOW\\(\\) WHERE id = \\$1 RETURNING id, username, password_hash, role, created_at, updated_at").
			WithArgs(uuid.MustParse(userID), userToUpdate.Username, userToUpdate.PasswordHash, userToUpdate.Role).
			WillReturnRows(rows)

		err := repo.Update(userToUpdate)
		require.NoError(t, err)
		assert.False(t, userToUpdate.UpdatedAt.IsZero())
	})

	t.Run("invalid user ID", func(t *testing.T) {
		invalidUser := &domain.User{ID: "invalid-uuid"}
		err := repo.Update(invalidUser)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID length")
	})

	t.Run("database error on update", func(t *testing.T) {
		mock.ExpectQuery("UPDATE users SET username = \\$2, password_hash = \\$3, role = \\$4, updated_at = NOW\\(\\) WHERE id = \\$1 RETURNING id, username, password_hash, role, created_at, updated_at").
			WithArgs(uuid.MustParse(userID), userToUpdate.Username, userToUpdate.PasswordHash, userToUpdate.Role).
			WillReturnError(errors.New("db error"))

		err := repo.Update(userToUpdate)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockQuerier := sqlc.New(db)
	repo := NewUserRepository(mockQuerier)

	userID := uuid.New().String()

	t.Run("successful user deletion", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
			WithArgs(uuid.MustParse(userID)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(userID)
		require.NoError(t, err)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		err := repo.Delete("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID length")
	})

	t.Run("database error on delete", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
			WithArgs(uuid.MustParse(userID)).
			WillReturnError(errors.New("db error"))

		err := repo.Delete(userID)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}
