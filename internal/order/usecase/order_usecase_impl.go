package usecase

import (
	"context"
	"gofiber-skeleton/internal/order/domain"
	"gofiber-skeleton/internal/order/infrastructure"
	"gofiber-skeleton/internal/infra/types"
)

type orderUsecase struct {
	orderRepo infrastructure.OrderRepository
}

func NewOrderUsecase(orderRepo infrastructure.OrderRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
	}
}

func (ou *orderUsecase) GetOrder(ctx context.Context, id types.UUIDv7) (*domain.Order, error) {
	return ou.orderRepo.GetOrder(ctx, id)
}

func (ou *orderUsecase) CreateOrder(ctx context.Context, order *domain.Order) error {
	return ou.orderRepo.CreateOrder(ctx, order)
}

func (ou *orderUsecase) UpdateOrder(ctx context.Context, order *domain.Order) error {
	return ou.orderRepo.UpdateOrder(ctx, order)
}

func (ou *orderUsecase) DeleteOrder(ctx context.Context, id types.UUIDv7) error {
	return ou.orderRepo.DeleteOrder(ctx, id)
}
