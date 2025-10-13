package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database/sqlc"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Post, error)
	ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type postRepository struct {
	queries Querier
}

type Querier interface {
	CreatePost(ctx context.Context, arg sqlc.CreatePostParams) (sqlc.Post, error)
	GetPostByID(ctx context.Context, id pgtype.UUID) (sqlc.Post, error)
	ListPosts(ctx context.Context, arg sqlc.ListPostsParams) ([]sqlc.Post, error)
	ListPostsByAuthor(ctx context.Context, arg sqlc.ListPostsByAuthorParams) ([]sqlc.Post, error)
	UpdatePost(ctx context.Context, arg sqlc.UpdatePostParams) (sqlc.Post, error)
	DeletePost(ctx context.Context, id pgtype.UUID) error
}

func NewPostRepository(queries Querier) PostRepository {
	return &postRepository{
		queries: queries,
	}
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) error {
	var idBytes [16]byte
	copy(idBytes[:], post.ID[:])

	var authorIDBytes [16]byte
	copy(authorIDBytes[:], post.AuthorID[:])

	params := sqlc.CreatePostParams{
		ID:        pgtype.UUID{Bytes: idBytes, Valid: true},
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  pgtype.UUID{Bytes: authorIDBytes, Valid: true},
		Status:    post.Status,
		CreatedAt: pgtype.Timestamptz{Time: post.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}
	_, err := r.queries.CreatePost(ctx, params)
	return err
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	dbPost, err := r.queries.GetPostByID(ctx, pgID)
	if err != nil {
		return nil, err
	}
	return r.dbPostToEntity(&dbPost), nil
}

func (r *postRepository) List(ctx context.Context, limit, offset int) ([]*entity.Post, error) {
	params := sqlc.ListPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	dbPosts, err := r.queries.ListPosts(ctx, params)
	if err != nil {
		return nil, err
	}

	posts := make([]*entity.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = r.dbPostToEntity(&dbPost)
	}
	return posts, nil
}

func (r *postRepository) ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*entity.Post, error) {
	var authorIDBytes [16]byte
	copy(authorIDBytes[:], authorID[:])

	pgAuthorID := pgtype.UUID{Bytes: authorIDBytes, Valid: true}
	params := sqlc.ListPostsByAuthorParams{
		AuthorID: pgAuthorID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}
	dbPosts, err := r.queries.ListPostsByAuthor(ctx, params)
	if err != nil {
		return nil, err
	}

	posts := make([]*entity.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = r.dbPostToEntity(&dbPost)
	}
	return posts, nil
}

func (r *postRepository) Update(ctx context.Context, post *entity.Post) error {
	var idBytes [16]byte
	copy(idBytes[:], post.ID[:])

	params := sqlc.UpdatePostParams{
		ID:        pgtype.UUID{Bytes: idBytes, Valid: true},
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}
	_, err := r.queries.UpdatePost(ctx, params)
	return err
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var idBytes [16]byte
	copy(idBytes[:], id[:])

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	return r.queries.DeletePost(ctx, pgID)
}

func (r *postRepository) dbPostToEntity(dbPost *sqlc.Post) *entity.Post {
	postID, _ := uuid.FromBytes(dbPost.ID.Bytes[:])
	authorID, _ := uuid.FromBytes(dbPost.AuthorID.Bytes[:])
	return &entity.Post{
		ID:        postID,
		Title:     dbPost.Title,
		Content:   dbPost.Content,
		AuthorID:  authorID,
		Status:    dbPost.Status,
		CreatedAt: dbPost.CreatedAt.Time,
		UpdatedAt: dbPost.UpdatedAt.Time,
	}
}