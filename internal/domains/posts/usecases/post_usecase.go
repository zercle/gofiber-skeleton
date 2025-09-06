package usecases

import (
	"context"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/models"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, authorID string, req models.CreatePostRequest) (*models.PostResponse, error)
	GetPost(ctx context.Context, id string) (*models.PostResponse, error)
	UpdatePost(ctx context.Context, id, authorID string, req models.UpdatePostRequest) (*models.PostResponse, error)
	DeletePost(ctx context.Context, id, authorID string) error
	ListPosts(ctx context.Context, page, pageSize int, publishedOnly bool) (*models.PostsListResponse, error)
	ListUserPosts(ctx context.Context, authorID string, page, pageSize int, publishedOnly bool) (*models.PostsListResponse, error)
	PublishPost(ctx context.Context, id, authorID string) (*models.PostResponse, error)
	UnpublishPost(ctx context.Context, id, authorID string) (*models.PostResponse, error)
}