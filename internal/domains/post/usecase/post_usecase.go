package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/repository"
)

var (
	ErrPostNotFound     = errors.New("post not found")
	ErrInvalidStatus    = errors.New("invalid post status")
	ErrUnauthorized     = errors.New("unauthorized to access this post")
)

type PostUsecase interface {
	Create(ctx context.Context, authorID uuid.UUID, req *entity.CreatePostRequest) (*entity.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Post, error)
	ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*entity.Post, error)
	Update(ctx context.Context, authorID uuid.UUID, id uuid.UUID, req *entity.UpdatePostRequest) (*entity.Post, error)
	Delete(ctx context.Context, authorID uuid.UUID, id uuid.UUID) error
}

type postUsecase struct {
	postRepo repository.PostRepository
}

func NewPostUsecase(postRepo repository.PostRepository) PostUsecase {
	return &postUsecase{
		postRepo: postRepo,
	}
}

func (u *postUsecase) Create(ctx context.Context, authorID uuid.UUID, req *entity.CreatePostRequest) (*entity.Post, error) {
	if !entity.IsValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	post := entity.NewPost(req.Title, req.Content, req.Status, authorID)

	if err := u.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (u *postUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (u *postUsecase) List(ctx context.Context, limit, offset int) ([]*entity.Post, error) {
	posts, err := u.postRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	return posts, nil
}

func (u *postUsecase) ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*entity.Post, error) {
	posts, err := u.postRepo.ListByAuthor(ctx, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by author: %w", err)
	}
	return posts, nil
}

func (u *postUsecase) Update(ctx context.Context, authorID uuid.UUID, id uuid.UUID, req *entity.UpdatePostRequest) (*entity.Post, error) {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}

	if post.AuthorID != authorID {
		return nil, ErrUnauthorized
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Status != nil {
		if !entity.IsValidStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		post.Status = *req.Status
	}

	post.UpdatedAt = time.Now()

	if err := u.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (u *postUsecase) Delete(ctx context.Context, authorID uuid.UUID, id uuid.UUID) error {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return ErrPostNotFound
	}

	if post.AuthorID != authorID {
		return ErrUnauthorized
	}

	if err := u.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}