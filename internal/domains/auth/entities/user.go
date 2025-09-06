package entities

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func NewUser(email, password, firstName, lastName string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id.String(),
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}