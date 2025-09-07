package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/entities"
)

// PostRepository defines the interface for post-related operations.
type PostRepository interface {
	Create(ctx context.Context, post *entities.Post) (*entities.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entities.Post, error)
	List(ctx context.Context, limit, offset int, isPublished *bool) ([]*entities.Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error)
	Update(ctx context.Context, post *entities.Post) (*entities.Post, error)
	Publish(ctx context.Context, id uuid.UUID) (*entities.Post, error)
	Unpublish(ctx context.Context, id uuid.UUID) (*entities.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
