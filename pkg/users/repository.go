package users

import (
	"fmt"

	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type userRepository struct {
	MainDbConn *gorm.DB
}

func InitUserRepository(mainDbConn *gorm.DB) domain.UserReposiroty {
	return &userRepository{
		MainDbConn: mainDbConn,
	}
}

func (r *userRepository) GetUser(userId string) (user models.User, err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.MainDbConn.Model(&models.User{})
	dbTx = dbTx.Where(models.User{Id: userId})
	err = dbTx.Take(&user).Error

	return
}

func (r *userRepository) GetUsers(criteria models.User) (users []models.User, err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.MainDbConn.Model(&models.User{})

	if len(criteria.Id) != 0 {
		dbTx = dbTx.Where(models.User{Id: criteria.Id})
	} else {
		if len(criteria.FullName) != 0 {
			dbTx = dbTx.Where("title LIKE ?", "%"+criteria.FullName+"%")
		}
	}

	err = dbTx.Find(&users).Error

	return
}

func (r *userRepository) CreateUser(user *models.User) (err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.MainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.User{})

	if err = dbTx.Create(user).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userRepository) EditUser(userId string, user models.User) (err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.MainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.User{})
	dbTx = dbTx.Where(models.User{Id: userId})

	if err = dbTx.Updates(user).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userRepository) DeleteUser(userId string) (err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.MainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.User{})
	dbTx = dbTx.Where(models.User{Id: userId})

	if err = dbTx.Delete(&models.User{}).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}
