package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type userReposiroty struct {
	DbConn *gorm.DB
}

func NewUserRepository(dbConn *gorm.DB) UserRepository {
	return &userReposiroty{
		DbConn: dbConn,
	}
}

func (r *userReposiroty) CreateUser(user *User) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&User{})

	if err = dbTx.Create(user).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error
	return
}

func (r *userReposiroty) UpdateUser(username string, user *User) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&User{})
	dbTx = dbTx.Where(User{Username: username})

	if err = dbTx.Updates(user).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userReposiroty) DeleteUser(username string) (err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&User{})
	dbTx = dbTx.Where(User{Username: username})

	if err = dbTx.Delete(&User{}).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userReposiroty) GetUser(username string) (user *User, err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Model(&User{})
	dbTx = dbTx.Where(User{Username: username})
	err = dbTx.Take(user).Error

	return
}

func (r *userReposiroty) GetUsers(fullname string) (users []User, err error) {
	if r.DbConn == nil {
		err = fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return
	}

	dbTx := r.DbConn.Model(&User{})

	if len(fullname) != 0 {
		dbTx = dbTx.Where("full_name LIKE ?", "%"+fullname+"%")
	}

	err = dbTx.Find(&users).Error

	return
}
