package users_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skelton/internal/datasources"
	"github.com/zercle/gofiber-skelton/internal/users"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetUserRepo(t *testing.T) {
	var mockUser models.User
	gofakeit.Struct(&mockUser)

	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDb,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "password", "full_name", "address", "create_at", "updated_at", "deleted_at"}).AddRow(mockUser.Id, mockUser.Password, mockUser.FullName, mockUser.Address, mockUser.CreatedAt, mockUser.UpdatedAt, mockUser.DeletedAt)

	query := "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL LIMIT 1"

	mock.ExpectQuery(query).WillReturnRows(rows)

	mockRepo := users.InitUserRepository(&datasources.Resources{MainDbConn: gdb})

	result, err := mockRepo.GetUser(mockUser.Id)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
