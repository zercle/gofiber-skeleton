package usecases

import (
	"context"
	"fmt"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/entities"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/models"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/repositories"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
)

type postUsecaseImpl struct {
	postRepo repositories.PostRepository
}

func NewPostUsecase(postRepo repositories.PostRepository) PostUsecase {
	return &postUsecaseImpl{
		postRepo: postRepo,
	}
}

func (u *postUsecaseImpl) CreatePost(ctx context.Context, authorID string, req models.CreatePostRequest) (*models.PostResponse, error) {
	post, err := entities.NewPost(req.Title, req.Content, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to create post entity: %w", err)
	}

	if err := u.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	response := models.NewPostResponse(post)
	return &response, nil
}

func (u *postUsecaseImpl) GetPost(ctx context.Context, id string) (*models.PostResponse, error) {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	response := models.NewPostResponse(post)
	return &response, nil
}

func (u *postUsecaseImpl) UpdatePost(ctx context.Context, id, authorID string, req models.UpdatePostRequest) (*models.PostResponse, error) {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post.AuthorID != authorID {
		return nil, types.ErrForbidden
	}

	post.Update(req.Title, req.Content)

	if err := u.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	response := models.NewPostResponse(post)
	return &response, nil
}

func (u *postUsecaseImpl) DeletePost(ctx context.Context, id, authorID string) error {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	if post.AuthorID != authorID {
		return types.ErrForbidden
	}

	if err := u.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (u *postUsecaseImpl) ListPosts(ctx context.Context, page, pageSize int, publishedOnly bool) (*models.PostsListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	posts, err := u.postRepo.List(ctx, pageSize, offset, publishedOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	total, err := u.postRepo.Count(ctx, publishedOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	response := models.NewPostsListResponse(posts, total, page, pageSize)
	return &response, nil
}

func (u *postUsecaseImpl) ListUserPosts(ctx context.Context, authorID string, page, pageSize int, publishedOnly bool) (*models.PostsListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	posts, err := u.postRepo.ListByAuthor(ctx, authorID, pageSize, offset, publishedOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list user posts: %w", err)
	}

	total, err := u.postRepo.CountByAuthor(ctx, authorID, publishedOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to count user posts: %w", err)
	}

	response := models.NewPostsListResponse(posts, total, page, pageSize)
	return &response, nil
}

func (u *postUsecaseImpl) PublishPost(ctx context.Context, id, authorID string) (*models.PostResponse, error) {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post.AuthorID != authorID {
		return nil, types.ErrForbidden
	}

	post.Publish()

	if err := u.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to publish post: %w", err)
	}

	response := models.NewPostResponse(post)
	return &response, nil
}

func (u *postUsecaseImpl) UnpublishPost(ctx context.Context, id, authorID string) (*models.PostResponse, error) {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post.AuthorID != authorID {
		return nil, types.ErrForbidden
	}

	post.Unpublish()

	if err := u.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to unpublish post: %w", err)
	}

	response := models.NewPostResponse(post)
	return &response, nil
}