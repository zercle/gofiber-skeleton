package infrastructure

import (
	"context"
	"gofiber-skeleton/internal/order/domain"
	"gofiber-skeleton/pkg/types"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (or *orderRepository) GetOrder(ctx context.Context, id types.UUIDv7) (*domain.Order, error) {
	var order domain.Order
	if err := or.db.WithContext(ctx).First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	if err := or.db.WithContext(ctx).Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	if err := or.db.WithContext(ctx).Save(order).Error; err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) DeleteOrder(ctx context.Context, id types.UUIDv7) error {
	if err := or.db.WithContext(ctx).Delete(&domain.Order{}, id).Error; err != nil {
		return err
	}
	return nil
}
