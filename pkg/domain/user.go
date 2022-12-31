package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type UserUsecase interface {
	GetUser(userId string) (user models.User, err error)
	GetUsers(criteria models.User) (users []models.User, err error)
	CreateUser(user *models.User) (err error)
	EditUser(userId string, user models.User) (err error)
	DeleteUser(userId string) (err error)
}

type UserReposiroty interface {
	GetUser(userId string) (user models.User, err error)
	GetUsers(criteria models.User) (users []models.User, err error)
	CreateUser(user *models.User) (err error)
	EditUser(userId string, user models.User) (err error)
	DeleteUser(userId string) (err error)
}
