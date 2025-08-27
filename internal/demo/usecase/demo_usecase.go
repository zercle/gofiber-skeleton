package usecase

import (
	"context"
	"fmt"

	demorepository "github.com/zercle/gofiber-skeleton/internal/demo/repository"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
)

// DemoUseCase defines the interface for demo-related business logic.
type DemoUseCase interface {
	PerformTransactionDemo(ctx context.Context) error
	GetJoinedDataDemo(ctx context.Context) ([]sqlc.GetOrdersWithItemsAndProductsRow, error)
}

// demoUseCase implements DemoUseCase.
type demoUseCase struct {
	demoRepo demorepository.DemoRepository
}

// NewDemoUseCase creates a new instance of DemoUseCase.
func NewDemoUseCase(demoRepo demorepository.DemoRepository) DemoUseCase {
	return &demoUseCase{
		demoRepo: demoRepo,
	}
}

// PerformTransactionDemo orchestrates the transaction demonstration.
func (uc *demoUseCase) PerformTransactionDemo(ctx context.Context) error {
	err := uc.demoRepo.PerformTransaction(ctx)
	if err != nil {
		return fmt.Errorf("usecase: failed to perform transaction demo: %w", err)
	}
	return nil
}

// GetJoinedDataDemo retrieves and processes the joined data.
func (uc *demoUseCase) GetJoinedDataDemo(ctx context.Context) ([]sqlc.GetOrdersWithItemsAndProductsRow, error) {
	data, err := uc.demoRepo.GetJoinedData(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to get joined data demo: %w", err)
	}
	return data, nil
}
