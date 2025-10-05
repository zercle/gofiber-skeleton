package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/db"
	"github.com/zercle/gofiber-skeleton/internal/post/entity"
)

//go:generate mockgen -source=postgres.go -destination=mocks/repository.go -package=mocks
type PostRepository interface {
	Create(ctx context.Context, p *entity.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Post, error)
	Update(ctx context.Context, p *entity.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type postgresPostRepository struct {
	queries *db.Queries
}

func NewPostgresPostRepository(queries *db.Queries) PostRepository {
	return &postgresPostRepository{queries: queries}
}

func (r *postgresPostRepository) Create(ctx context.Context, p *entity.Post) error {
	result, err := r.queries.CreatePost(ctx, db.CreatePostParams{
		UserID:   p.UserID,
		ThreadID: p.ThreadID,
		Content:  p.Content,
	})
	if err != nil {
		return err
	}
	// Update the entity with generated values
	p.ID = result.ID
	p.CreatedAt = result.CreatedAt.Time
	p.UpdatedAt = result.UpdatedAt.Time
	return nil
}

func (r *postgresPostRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	post, err := r.queries.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		ID:        post.ID,
		ThreadID:  post.ThreadID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Time,
		UpdatedAt: post.UpdatedAt.Time,
	}, nil
}

func (r *postgresPostRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Post, error) {
	posts, err := r.queries.ListPostsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var postEntities []*entity.Post
	for _, post := range posts {
		postEntities = append(postEntities, &entity.Post{
			ID:        post.ID,
			ThreadID:  post.ThreadID,
			UserID:    post.UserID,
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}
	return postEntities, nil
}

func (r *postgresPostRepository) Update(ctx context.Context, p *entity.Post) error {
	result, err := r.queries.UpdatePost(ctx, db.UpdatePostParams{
		ID:      p.ID,
		Content: p.Content,
	})
	if err != nil {
		return err
	}
	// Update the entity with new values
	p.UpdatedAt = result.UpdatedAt.Time
	return nil
}

func (r *postgresPostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePost(ctx, id)
}
