package domain

import (
	product_domain "gofiber-skeleton/internal/product/domain"
	user_domain "gofiber-skeleton/internal/user/domain"
)

type Order struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	User      user_domain.User    `gorm:"foreignKey:UserID"`
	ProductID uint    `gorm:"not null"`
	Product   product_domain.Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
}
