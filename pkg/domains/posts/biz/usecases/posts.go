package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/pkg/domains/posts/entities"
	"github.com/zercle/gofiber-skeleton/pkg/domains/posts/models"
	"github.com/zercle/gofiber-skeleton/pkg/domains/posts/repositories"
)

// PostUseCase defines the interface for post-related use cases.
type PostUseCase interface {
	Create(ctx context.Context, userID uuid.UUID, req *models.CreatePostRequest) (*models.PostResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.PostWithAuthorResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.PostWithAuthorResponse, error)
	List(ctx context.Context, limit, offset int, isPublished *bool) (*models.PostListResponse, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.PostListResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdatePostRequest) (*models.PostResponse, error)
	Publish(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.PostResponse, error)
	Unpublish(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.PostResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type postUseCase struct {
	postRepo repositories.PostRepository
}

// NewPostUseCase creates a new instance of PostUseCase.
func NewPostUseCase(postRepo repositories.PostRepository) PostUseCase {
	return &postUseCase{
		postRepo: postRepo,
	}
}

func (p *postUseCase) Create(ctx context.Context, userID uuid.UUID, req *models.CreatePostRequest) (*models.PostResponse, error) {
	now := time.Now()
	post := &entities.Post{
		UserID:      userID,
		Title:       req.Title,
		Content:     req.Content,
		Slug:        p.generateSlug(req.Title),
		IsPublished: req.IsPublished,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if req.IsPublished {
		post.PublishedAt = &now
	}

	createdPost, err := p.postRepo.Create(ctx, post)
	if err != nil {
		return nil, errors.New("failed to create post")
	}

	return p.toPostResponse(createdPost), nil
}

func (p *postUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.PostWithAuthorResponse, error) {
	post, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return p.toPostWithAuthorResponse(post), nil
}

func (p *postUseCase) GetBySlug(ctx context.Context, slug string) (*models.PostWithAuthorResponse, error) {
	post, err := p.postRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return p.toPostWithAuthorResponse(post), nil
}

func (p *postUseCase) List(ctx context.Context, limit, offset int, isPublished *bool) (*models.PostListResponse, error) {
	posts, err := p.postRepo.List(ctx, limit, offset, isPublished)
	if err != nil {
		return nil, errors.New("failed to fetch posts")
	}

	responses := make([]models.PostWithAuthorResponse, len(posts))
	for i, post := range posts {
		responses[i] = *p.toPostWithAuthorResponse(post)
	}

	return &models.PostListResponse{
		Posts:  responses,
		Total:  len(responses), // This should be replaced with actual count from DB
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (p *postUseCase) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.PostListResponse, error) {
	posts, err := p.postRepo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch user posts")
	}

	responses := make([]models.PostWithAuthorResponse, len(posts))
	for i, post := range posts {
		responses[i] = *p.toPostWithAuthorResponse(post)
	}

	return &models.PostListResponse{
		Posts:  responses,
		Total:  len(responses), // This should be replaced with actual count from DB
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (p *postUseCase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdatePostRequest) (*models.PostResponse, error) {
	post, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized: you can only update your own posts")
	}

	if req.Title != nil {
		post.Title = *req.Title
		post.Slug = p.generateSlug(*req.Title)
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.IsPublished != nil {
		post.IsPublished = *req.IsPublished
		if *req.IsPublished && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		} else if !*req.IsPublished {
			post.PublishedAt = nil
		}
	}

	post.UpdatedAt = time.Now()

	updatedPost, err := p.postRepo.Update(ctx, post)
	if err != nil {
		return nil, errors.New("failed to update post")
	}

	return p.toPostResponse(updatedPost), nil
}

func (p *postUseCase) Publish(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.PostResponse, error) {
	post, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized: you can only publish your own posts")
	}

	publishedPost, err := p.postRepo.Publish(ctx, id)
	if err != nil {
		return nil, errors.New("failed to publish post")
	}

	return p.toPostResponse(publishedPost), nil
}

func (p *postUseCase) Unpublish(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.PostResponse, error) {
	post, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized: you can only unpublish your own posts")
	}

	unpublishedPost, err := p.postRepo.Unpublish(ctx, id)
	if err != nil {
		return nil, errors.New("failed to unpublish post")
	}

	return p.toPostResponse(unpublishedPost), nil
}

func (p *postUseCase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	post, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("post not found")
	}

	if post.UserID != userID {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	if err := p.postRepo.Delete(ctx, id); err != nil {
		return errors.New("failed to delete post")
	}

	return nil
}

func (p *postUseCase) generateSlug(title string) string {
	// Simple slug generation - replace spaces with dashes and make lowercase
	// In production, you might want a more robust slug generation
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Add timestamp to ensure uniqueness
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", slug, timestamp)
}

func (p *postUseCase) toPostResponse(post *entities.Post) *models.PostResponse {
	return &models.PostResponse{
		ID:          post.ID,
		UserID:      post.UserID,
		Title:       post.Title,
		Content:     post.Content,
		Slug:        post.Slug,
		IsPublished: post.IsPublished,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}

func (p *postUseCase) toPostWithAuthorResponse(post *entities.Post) *models.PostWithAuthorResponse {
	// For now, return without author info since we need to implement the repo method
	// This would need to be updated when the repo returns PostWithAuthor entity
	return &models.PostWithAuthorResponse{
		ID:              post.ID,
		UserID:          post.UserID,
		Title:           post.Title,
		Content:         post.Content,
		Slug:            post.Slug,
		IsPublished:     post.IsPublished,
		PublishedAt:     post.PublishedAt,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		AuthorEmail:     "", // These would come from joined user data
		AuthorFirstName: "",
		AuthorLastName:  "",
		AuthorFullName:  "",
	}
}
