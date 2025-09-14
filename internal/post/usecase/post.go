package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/post/repository"
)

//go:generate mockgen -source=post.go -destination=mocks/usecase.go -package=mocks
type PostUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, threadID uuid.UUID, content string) (*entity.Post, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Post, error)
	Update(ctx context.Context, userID uuid.UUID, p *entity.Post) (*entity.Post, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type postUsecase struct {
	postRepo repository.PostRepository
}

func NewPostUsecase(postRepo repository.PostRepository) PostUsecase {
	return &postUsecase{postRepo: postRepo}
}

func (uc *postUsecase) Create(ctx context.Context, userID uuid.UUID, threadID uuid.UUID, content string) (*entity.Post, error) {
	var err error
	
	post := &entity.Post{
		UserID:    userID,
		ThreadID:  threadID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if post.ID, err = uuid.NewV7(); err != nil {
		return nil, fmt.Errorf("failed to generate post ID: %w", err)
	}

	if err := uc.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (uc *postUsecase) Get(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	return post, nil
}

func (uc *postUsecase) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Post, error) {
	posts, err := uc.postRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by user: %w", err)
	}
	return posts, nil
}

func (uc *postUsecase) Update(ctx context.Context, userID uuid.UUID, p *entity.Post) (*entity.Post, error) {
	existingPost, err := uc.postRepo.GetByID(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post for update: %w", err)
	}

	if existingPost.UserID != userID {
		return nil, fmt.Errorf("user is not the owner of this post")
	}

	existingPost.Content = p.Content
	existingPost.UpdatedAt = time.Now()

	if err := uc.postRepo.Update(ctx, existingPost); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return existingPost, nil
}

func (uc *postUsecase) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	existingPost, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get post for delete: %w", err)
	}

	if existingPost.UserID != userID {
		return fmt.Errorf("user is not the owner of this post")
	}

	if err := uc.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}