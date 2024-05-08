package models

import (
	"database/sql"

	"github.com/oklog/ulid/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"size:32;primaryKey"`
	Password  string         `json:"password" gorm:"size:64"`
	FullName  string         `json:"full_name" gorm:"size:127;index"`
	Address   string         `json:"address" gorm:"type:text"`
	CreatedAt sql.NullTime   `json:"createdAt,omitempty" gorm:"autoCreateTime;index"`
	UpdatedAt sql.NullTime   `json:"updatedAt,omitempty" gorm:"autoUpdateTime;index"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.ID) == 0 {
		u.ID = ulid.Make().String()
	}
	return
}

type UserResponse struct {
	helpers.ResponseForm
}
