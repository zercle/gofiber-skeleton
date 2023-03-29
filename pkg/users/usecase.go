package users

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type userUsecase struct {
	userRepo domain.UserReposiroty
}

func NewUserUsecase(userRepo domain.UserReposiroty) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetUser(userId string) (user models.User, err error) {
	return u.userRepo.GetUser(userId)
}

func (u *userUsecase) GetUsers(criteria models.User) (users []models.User, err error) {
	return u.userRepo.GetUsers(criteria)
}

func (u *userUsecase) CreateUser(user *models.User) (err error) {
	return u.userRepo.CreateUser(user)
}

func (u *userUsecase) EditUser(userId string, user models.User) (err error) {
	return u.userRepo.EditUser(userId, user)
}

func (u *userUsecase) DeleteUser(userId string) (err error) {
	return u.userRepo.DeleteUser(userId)
}
