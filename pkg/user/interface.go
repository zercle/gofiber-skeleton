package user

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
