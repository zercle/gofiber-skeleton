package users_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/users"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetUserRepo(t *testing.T) {
	var mockUser models.User
	gofakeit.Struct(&mockUser)

	mockUser.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	mockUser.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	mockUser.DeletedAt = gorm.DeletedAt(sql.NullTime{Valid: false})

	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDb,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "password", "full_name", "address", "created_at", "updated_at", "deleted_at"}).AddRow(mockUser.Id, mockUser.Password, mockUser.FullName, mockUser.Address, mockUser.CreatedAt, mockUser.UpdatedAt, mockUser.DeletedAt)

	queryRegexp := "^SELECT (.+) FROM `users` (.+)$"

	mock.ExpectQuery(queryRegexp).WillReturnRows(rows)

	t.Run("success", func(t *testing.T) {
		mockRepo := users.InitUserRepository(gdb)

		result, err := mockRepo.GetUser(mockUser.Id)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, mockUser, result)
	})

	t.Run("fail-db-conn", func(t *testing.T) {
		mockRepo := users.InitUserRepository(nil)

		result, err := mockRepo.GetUser(mockUser.Id)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
