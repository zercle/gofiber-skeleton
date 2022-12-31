package users

import (
	"fmt"

	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/internal/datasources"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type userRepository struct {
	*datasources.Resources
}

func InitUserRepository(resources *datasources.Resources) domain.UserReposiroty {
	return &userRepository{
		Resources: resources,
	}
}

func (r *userRepository) GetUser(userId string) (user models.User, err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}
	return
}

func (r *userRepository) GetUsers(criteria models.User) (users []models.User, err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}
	return
}

func (r *userRepository) CreateUser(user *models.User) (err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}
	return
}

func (r *userRepository) EditUser(userId string, user models.User) (err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}
	return
}

func (r *userRepository) DeleteUser(userId string) (err error) {
	if r.MainDbConn == nil {
		err = fmt.Errorf("%s \nErr: %+v", helpers.WhereAmI(), "database has gone away.")
		return
	}
	return
}
