package users

import (
	"fmt"

	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type userRepository struct {
	mainDbConn *gorm.DB
}

func NewUserRepository(mainDbConn *gorm.DB) domain.UserReposiroty {
	return &userRepository{
		mainDbConn: mainDbConn,
	}
}

func (r *userRepository) GetUser(userID string) (user models.User, err error) {
	if r.mainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.mainDbConn.Model(&models.User{})
	dbTx = dbTx.Where(models.User{ID: userID})
	err = dbTx.Take(&user).Error

	return
}

func (r *userRepository) GetUsers(criteria models.User) (users []models.User, err error) {
	if r.mainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.mainDbConn.Model(&models.User{})

	if len(criteria.ID) != 0 {
		dbTx = dbTx.Where(models.User{ID: criteria.ID})
	} else {
		if len(criteria.FullName) != 0 {
			dbTx = dbTx.Where("title LIKE ?", "%"+criteria.FullName+"%")
		}
	}

	err = dbTx.Find(&users).Error

	return
}

func (r *userRepository) CreateUser(user *models.User) (err error) {
	if r.mainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.mainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.User{})

	if err = dbTx.Create(user).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userRepository) EditUser(userID string, user models.User) (err error) {
	if r.mainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.mainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.User{})
	dbTx = dbTx.Where(models.User{ID: userID})

	if err = dbTx.Updates(user).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userRepository) DeleteUser(userID string) (err error) {
	if r.mainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}

	dbTx := r.mainDbConn.Begin()
	defer dbTx.Rollback()

	dbTx = dbTx.Model(&models.User{})
	dbTx = dbTx.Where(models.User{ID: userID})

	if err = dbTx.Delete(&models.User{}).Error; err != nil {
		return
	}

	err = dbTx.Commit().Error

	return
}
