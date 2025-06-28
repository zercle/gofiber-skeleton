package usecase

import (
	"context"
	"gofiber-skeleton/internal/order/domain"
	"gofiber-skeleton/internal/order/infrastructure"
)

type orderUsecase struct {
	orderRepo infrastructure.OrderRepository
}

func NewOrderUsecase(orderRepo infrastructure.OrderRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
	}
}

func (ou *orderUsecase) GetOrder(ctx context.Context, id uint) (*domain.Order, error) {
	return ou.orderRepo.GetOrder(ctx, id)
}

func (ou *orderUsecase) CreateOrder(ctx context.Context, order *domain.Order) error {
	return ou.orderRepo.CreateOrder(ctx, order)
}

func (ou *orderUsecase) UpdateOrder(ctx context.Context, order *domain.Order) error {
	return ou.orderRepo.UpdateOrder(ctx, order)
}

func (ou *orderUsecase) DeleteOrder(ctx context.Context, id uint) error {
	return ou.orderRepo.DeleteOrder(ctx, id)
}
