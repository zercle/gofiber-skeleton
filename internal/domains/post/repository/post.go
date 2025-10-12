package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
)

//go:generate mockgen -source=post.go -destination=../mocks/post_repository_mock.go -package=mocks

type PostRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, post *entity.DomainPost) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.DomainPost, error)
	Update(ctx context.Context, post *entity.DomainPost) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User-specific operations
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.DomainPost, error)
	GetPublishedPostsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.DomainPost, error)
	GetDraftPostsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.DomainPost, error)

	// Public operations with author info
	GetAllPosts(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error)
	GetPostsWithAuthor(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error)
	SearchPosts(ctx context.Context, search string, status entity.PostStatus, limit, offset int) ([]*entity.PostWithAuthor, error)

	// Statistics and aggregation
	GetUserPostStats(ctx context.Context, userID uuid.UUID) (*entity.PostStats, error)
	CountPostsByStatus(ctx context.Context) ([]*entity.PostStatusCount, error)

	// Authorization checks
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	IsOwner(ctx context.Context, postID, userID uuid.UUID) (bool, error)
}
