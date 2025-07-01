//go:generate mockgen -source=order_usecase.go -destination=../mocks/mock_order_usecase.go -package=mocks
package usecase

import (
	"context"
	"gofiber-skeleton/internal/order/domain"
	"gofiber-skeleton/pkg/types"
)

type OrderUsecase interface {
	GetOrder(ctx context.Context, id types.UUIDv7) (*domain.Order, error)
	CreateOrder(ctx context.Context, order *domain.Order) error
	UpdateOrder(ctx context.Context, order *domain.Order) error
	DeleteOrder(ctx context.Context, id types.UUIDv7) error
}