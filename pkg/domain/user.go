package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type UserUsecase interface {
	GetUser(userID string) (user models.User, err error)
	GetUsers(criteria models.User) (users []models.User, err error)
	CreateUser(user *models.User) (err error)
	EditUser(userID string, user models.User) (err error)
	DeleteUser(userID string) (err error)
}

type UserReposiroty interface {
	GetUser(userID string) (user models.User, err error)
	GetUsers(criteria models.User) (users []models.User, err error)
	CreateUser(user *models.User) (err error)
	EditUser(userID string, user models.User) (err error)
	DeleteUser(userID string) (err error)
}
