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

func (u *userUsecase) GetUser(userID string) (user models.User, err error) {
	return u.userRepo.GetUser(userID)
}

func (u *userUsecase) GetUsers(criteria models.User) (users []models.User, err error) {
	return u.userRepo.GetUsers(criteria)
}

func (u *userUsecase) CreateUser(user *models.User) (err error) {
	return u.userRepo.CreateUser(user)
}

func (u *userUsecase) EditUser(userID string, user models.User) (err error) {
	return u.userRepo.EditUser(userID, user)
}

func (u *userUsecase) DeleteUser(userID string) (err error) {
	return u.userRepo.DeleteUser(userID)
}
