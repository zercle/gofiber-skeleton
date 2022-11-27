package user

import (
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type userService struct {
	userRepository UserRepository
}

func NewUserService(r UserRepository) UserService {
	return &userService{
		userRepository: r,
	}
}

func (s *userService) CreateUser(user *User) (err error) {
	if len(user.Username) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: username")
		return
	}
	if len(user.RawPassword) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: raw_password")
		return
	}
	if user.RawPassword != user.Password {
		err = fiber.NewError(fiber.StatusBadRequest, "need: password miss match")
		return
	}
	if len(user.FullName) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: full_name")
		return
	}

	if user.Password, err = helpers.HashPasswordString(user.RawPassword); err != nil {
		return
	}

	user.RawPassword = ""

	return s.userRepository.CreateUser(user)
}

func (s *userService) ChangePassword(username, oldPassword, newPassword, confirmPassword string) (err error) {
	if len(username) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: username")
		return
	}

	if len(newPassword) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: raw_password")
		return
	}
	if newPassword != confirmPassword {
		err = fiber.NewError(fiber.StatusBadRequest, "need: new_password miss match")
		return
	}

	currentUser, err := s.userRepository.GetUser(username)
	if err != nil {
		return
	}

	if len(oldPassword) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: old_password")
		return
	}
	if err = helpers.CheckPasswordHashString(oldPassword, currentUser.Password); err != nil {
		err = fiber.NewError(fiber.StatusBadRequest, "need: old_password miss match")
		return
	}

	user := new(User)

	if user.Password, err = helpers.HashPasswordString(newPassword); err != nil {
		return
	}

	user.RawPassword = ""

	return s.userRepository.UpdateUser(username, user)
}

func (s *userService) UpdateUser(username string, user *User) (err error) {
	if len(username) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: username")
		return
	}

	currentUser, err := s.userRepository.GetUser(username)
	if err != nil {
		return
	}

	if len(user.RawPassword) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: raw_password")
		return
	}
	if err = helpers.CheckPasswordHashString(user.RawPassword, currentUser.Password); err != nil {
		err = fiber.NewError(fiber.StatusBadRequest, "need: password miss match")
		return
	}

	user.Password = ""
	user.RawPassword = ""

	return s.userRepository.UpdateUser(username, user)
}

func (s *userService) DeleteUser(username string) (err error) {
	if len(username) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: username")
		return
	}
	return s.userRepository.DeleteUser(username)
}

func (s *userService) GetUser(username string) (user *User, err error) {
	if len(username) == 0 {
		err = fiber.NewError(fiber.StatusBadRequest, "need: username")
		return
	}
	return s.userRepository.GetUser(username)
}

func (s *userService) GetUsers(fullname string) (users []User, err error) {
	return s.userRepository.GetUsers(fullname)
}
