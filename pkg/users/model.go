package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username    string         `json:"username" gorm:"primaryKey"`
	Password    string         `json:"password" gorm:""`
	RawPassword string         `json:"raw_password" gorm:"-"`
	FullName    string         `json:"full_name" gorm:"size:127;index"`
	CreatedAt   *time.Time     `json:"createdAt" gorm:"autoCreateTime;index"`
	UpdatedAt   *time.Time     `json:"updatedAt" gorm:"autoUpdateTime;index"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type UserRepository interface {
	GetUsers(fullname string) (*[]User, error)
	GetUser(username string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(username string, user *User) error
	DeleteUser(username string) error
}

type UserService interface {
	GetUsers(fullname string) (*[]User, error)
	GetUser(username string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(username string, user *User) error
	DeleteUser(username string) error
	ChangePassword(username string, oldPassword string, newPassword string, confirmPassword string) error
}
