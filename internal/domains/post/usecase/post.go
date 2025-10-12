package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/repository"
)

type PostUsecase interface {
	// Post management
	CreatePost(ctx context.Context, userID uuid.UUID, req *entity.CreatePostRequest) (*entity.DomainPost, error)
	GetPost(ctx context.Context, postID uuid.UUID) (*entity.DomainPost, error)
	UpdatePost(ctx context.Context, postID, userID uuid.UUID, req *entity.UpdatePostRequest) (*entity.DomainPost, error)
	DeletePost(ctx context.Context, postID, userID uuid.UUID) error

	// User posts
	GetUserPosts(ctx context.Context, userID uuid.UUID, status *entity.PostStatus) ([]*entity.DomainPost, error)
	GetUserPostStats(ctx context.Context, userID uuid.UUID) (*entity.PostStats, error)

	// Public posts
	GetAllPosts(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error)
	GetPublishedPosts(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error)
	SearchPosts(ctx context.Context, search string, limit, offset int) ([]*entity.PostWithAuthor, error)

	// Post operations
	PublishPost(ctx context.Context, postID, userID uuid.UUID) error
	ArchivePost(ctx context.Context, postID, userID uuid.UUID) error
	UnpublishPost(ctx context.Context, postID, userID uuid.UUID) error
}

type postUsecase struct {
	postRepo repository.PostRepository
}

func NewPostUsecase(injector do.Injector) (PostUsecase, error) {
	postRepo := do.MustInvoke[repository.PostRepository](injector)
	return &postUsecase{
		postRepo: postRepo,
	}, nil
}

// CreatePost creates a new post
func (u *postUsecase) CreatePost(ctx context.Context, userID uuid.UUID, req *entity.CreatePostRequest) (*entity.DomainPost, error) {
	// Validate request
	if req.Title == "" {
		return nil, entity.ErrPostTitleRequired
	}
	if req.Content == "" {
		return nil, entity.ErrPostContentRequired
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = entity.PostStatusDraft
	}

	// Validate status
	if status != entity.PostStatusDraft && status != entity.PostStatusPublished && status != entity.PostStatusArchived {
		return nil, entity.ErrInvalidPostStatus
	}

	post := entity.NewPost(req.Title, req.Content, userID)
	post.Status = status

	if err := post.Validate(); err != nil {
		return nil, err
	}

	if err := u.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// GetPost retrieves a post by ID
func (u *postUsecase) GetPost(ctx context.Context, postID uuid.UUID) (*entity.DomainPost, error) {
	return u.postRepo.GetByID(ctx, postID)
}

// UpdatePost updates an existing post
func (u *postUsecase) UpdatePost(ctx context.Context, postID, userID uuid.UUID, req *entity.UpdatePostRequest) (*entity.DomainPost, error) {
	// Check if post exists and user is owner
	exists, err := u.postRepo.IsOwner(ctx, postID, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entity.ErrPostNotFound
	}

	// Get existing post
	post, err := u.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != nil {
		if *req.Title == "" {
			return nil, entity.ErrPostTitleRequired
		}
		post.Title = *req.Title
	}
	if req.Content != nil {
		if *req.Content == "" {
			return nil, entity.ErrPostContentRequired
		}
		post.Content = *req.Content
	}
	if req.Status != nil {
		if *req.Status != entity.PostStatusDraft && *req.Status != entity.PostStatusPublished && *req.Status != entity.PostStatusArchived {
			return nil, entity.ErrInvalidPostStatus
		}
		post.Status = *req.Status
	}

	post.UpdatedAt = time.Now()

	if err := post.Validate(); err != nil {
		return nil, err
	}

	if err := u.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// DeletePost deletes a post
func (u *postUsecase) DeletePost(ctx context.Context, postID, userID uuid.UUID) error {
	// Check if post exists and user is owner
	exists, err := u.postRepo.IsOwner(ctx, postID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return entity.ErrPostNotFound
	}

	return u.postRepo.Delete(ctx, postID)
}

// GetUserPosts retrieves posts for a specific user
func (u *postUsecase) GetUserPosts(ctx context.Context, userID uuid.UUID, status *entity.PostStatus) ([]*entity.DomainPost, error) {
	if status != nil {
		switch *status {
		case entity.PostStatusPublished:
			return u.postRepo.GetPublishedPostsByUserID(ctx, userID)
		case entity.PostStatusDraft:
			return u.postRepo.GetDraftPostsByUserID(ctx, userID)
		default:
			return nil, entity.ErrInvalidPostStatus
		}
	}
	return u.postRepo.GetByUserID(ctx, userID)
}

// GetUserPostStats retrieves post statistics for a user
func (u *postUsecase) GetUserPostStats(ctx context.Context, userID uuid.UUID) (*entity.PostStats, error) {
	return u.postRepo.GetUserPostStats(ctx, userID)
}

// GetAllPosts retrieves all posts with author info
func (u *postUsecase) GetAllPosts(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return u.postRepo.GetAllPosts(ctx, limit, offset)
}

// GetPublishedPosts retrieves published posts with author info
func (u *postUsecase) GetPublishedPosts(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return u.postRepo.GetPostsWithAuthor(ctx, limit, offset)
}

// SearchPosts searches published posts
func (u *postUsecase) SearchPosts(ctx context.Context, search string, limit, offset int) ([]*entity.PostWithAuthor, error) {
	if search == "" {
		return nil, errors.New("search query is required")
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return u.postRepo.SearchPosts(ctx, search, entity.PostStatusPublished, limit, offset)
}

// PublishPost publishes a post
func (u *postUsecase) PublishPost(ctx context.Context, postID, userID uuid.UUID) error {
	return u.updatePostStatus(ctx, postID, userID, entity.PostStatusPublished)
}

// ArchivePost archives a post
func (u *postUsecase) ArchivePost(ctx context.Context, postID, userID uuid.UUID) error {
	return u.updatePostStatus(ctx, postID, userID, entity.PostStatusArchived)
}

// UnpublishPost unpublishes a post (sets to draft)
func (u *postUsecase) UnpublishPost(ctx context.Context, postID, userID uuid.UUID) error {
	return u.updatePostStatus(ctx, postID, userID, entity.PostStatusDraft)
}

// updatePostStatus is a helper method to update post status
func (u *postUsecase) updatePostStatus(ctx context.Context, postID, userID uuid.UUID, status entity.PostStatus) error {
	// Check if post exists and user is owner
	exists, err := u.postRepo.IsOwner(ctx, postID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return entity.ErrPostNotFound
	}

	// Get existing post
	post, err := u.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Update status
	post.Status = status
	post.UpdatedAt = time.Now()

	return u.postRepo.Update(ctx, post)
}
