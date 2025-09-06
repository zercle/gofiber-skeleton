package repositories

import (
	"context"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/entities"
)

type PostRepository interface {
	Create(ctx context.Context, post *entities.Post) error
	GetByID(ctx context.Context, id string) (*entities.Post, error)
	Update(ctx context.Context, post *entities.Post) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int, publishedOnly bool) ([]*entities.Post, error)
	ListByAuthor(ctx context.Context, authorID string, limit, offset int, publishedOnly bool) ([]*entities.Post, error)
	Count(ctx context.Context, publishedOnly bool) (int64, error)
	CountByAuthor(ctx context.Context, authorID string, publishedOnly bool) (int64, error)
}