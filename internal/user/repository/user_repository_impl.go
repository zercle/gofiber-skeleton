package repository

import (
	"context"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/infra/types"

	"gorm.io/gorm"
)

var _ UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := ur.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetUser(ctx context.Context, id types.UUIDv7) (*domain.User, error) {
	var user domain.User
	if err := ur.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if err := ur.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := ur.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, id types.UUIDv7) error {
	if err := ur.db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
