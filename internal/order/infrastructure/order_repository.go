//go:generate mockgen -source=order_repository.go -destination=../mocks/mock_order_repository.go -package=mocks
package infrastructure

import (
	"context"
	"gofiber-skeleton/internal/order/domain"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, id uint) (*domain.Order, error)
	CreateOrder(ctx context.Context, order *domain.Order) error
	UpdateOrder(ctx context.Context, order *domain.Order) error
	DeleteOrder(ctx context.Context, id uint) error
}
