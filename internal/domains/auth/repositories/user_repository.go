package repositories

import (
	"context"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	Count(ctx context.Context) (int64, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}