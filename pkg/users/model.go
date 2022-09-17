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
	CreatedAt   *time.Time      `json:"createdAt" gorm:"autoCreateTime;index"`
	UpdatedAt   *time.Time      `json:"updatedAt" gorm:"autoUpdateTime;index"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}


