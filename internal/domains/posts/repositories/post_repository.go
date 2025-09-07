package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/entities"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	sqldb "github.com/zercle/gofiber-skeleton/internal/infrastructure/database/queries"
)

type postRepository struct {
	q *sqldb.Queries
}

func NewPostRepository(db *database.Database) PostRepository {
	return &postRepository{q: sqldb.New(db.Pool)}
}

func (r *postRepository) Create(ctx context.Context, post *entities.Post) (*entities.Post, error) {
	params := sqldb.CreatePostParams{
		UserID:  pgtype.UUID{Bytes: post.UserID, Valid: true},
		Title:   post.Title,
		Content: pgtype.Text{String: post.Content, Valid: true},
		Slug:    post.Slug,
	}
	createdPost, err := r.q.CreatePost(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	return r.mapSQLCPostToEntity(createdPost), nil
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	sqlcID := pgtype.UUID{Bytes: id, Valid: true}
	post, err := r.q.GetPostByID(ctx, sqlcID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Post not found
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	// Convert PostWithAuthor to Post
	postWithAuthor := r.mapSQLCGetPostByIDRowToEntity(post)
	return &postWithAuthor.Post, nil
}

func (r *postRepository) GetBySlug(ctx context.Context, slug string) (*entities.Post, error) {
	post, err := r.q.GetPostBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Post not found
		}
		return nil, fmt.Errorf("failed to get post by slug: %w", err)
	}
	// Convert PostWithAuthor to Post
	postWithAuthor := r.mapSQLCGetPostBySlugRowToEntity(post)
	return &postWithAuthor.Post, nil
}

func (r *postRepository) List(ctx context.Context, limit, offset int, isPublished *bool) ([]*entities.Post, error) {
	params := sqldb.ListPostsParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Column3: isPublished != nil && *isPublished,
	}
	posts, err := r.q.ListPosts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	var postEntities []*entities.Post
	for _, post := range posts {
		// Convert PostWithAuthor to Post
		postWithAuthor := r.mapSQLCListPostsRowToEntity(post)
		postEntities = append(postEntities, &postWithAuthor.Post)
	}
	return postEntities, nil
}

func (r *postRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error) {
	params := sqldb.ListPostsByUserParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	posts, err := r.q.ListPostsByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by user: %w", err)
	}

	var postEntities []*entities.Post
	for _, post := range posts {
		postEntities = append(postEntities, r.mapSQLCPostToEntity(post))
	}
	return postEntities, nil
}

func (r *postRepository) Update(ctx context.Context, post *entities.Post) (*entities.Post, error) {
	params := sqldb.UpdatePostParams{
		ID:      pgtype.UUID{Bytes: post.ID, Valid: true},
		Title:   post.Title,
		Content: pgtype.Text{String: post.Content, Valid: true},
		Slug:    post.Slug,
	}
	updatedPost, err := r.q.UpdatePost(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	return r.mapSQLCPostToEntity(updatedPost), nil
}

func (r *postRepository) Publish(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	sqlcID := pgtype.UUID{Bytes: id, Valid: true}
	publishedPost, err := r.q.PublishPost(ctx, sqlcID)
	if err != nil {
		return nil, fmt.Errorf("failed to publish post: %w", err)
	}
	return r.mapSQLCPostToEntity(publishedPost), nil
}

func (r *postRepository) Unpublish(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	sqlcID := pgtype.UUID{Bytes: id, Valid: true}
	unpublishedPost, err := r.q.UnpublishPost(ctx, sqlcID)
	if err != nil {
		return nil, fmt.Errorf("failed to unpublish post: %w", err)
	}
	return r.mapSQLCPostToEntity(unpublishedPost), nil
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sqlcID := pgtype.UUID{Bytes: id, Valid: true}
	err := r.q.DeletePost(ctx, sqlcID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (r *postRepository) mapSQLCPostToEntity(sqlcPost sqldb.Posts) *entities.Post {
	return &entities.Post{
		ID:          uuid.UUID(sqlcPost.ID.Bytes),
		UserID:      uuid.UUID(sqlcPost.UserID.Bytes),
		Title:       sqlcPost.Title,
		Content:     sqlcPost.Content.String,
		Slug:        sqlcPost.Slug,
		IsPublished: sqlcPost.IsPublished.Bool,
		PublishedAt: r.pgtypeTimestampToTimePtr(sqlcPost.PublishedAt),
		CreatedAt:   sqlcPost.CreatedAt.Time,
		UpdatedAt:   sqlcPost.UpdatedAt.Time,
	}
}

func (r *postRepository) mapSQLCGetPostByIDRowToEntity(sqlcPost sqldb.GetPostByIDRow) *entities.PostWithAuthor {
	return &entities.PostWithAuthor{
		Post: entities.Post{
			ID:          uuid.UUID(sqlcPost.ID.Bytes),
			UserID:      uuid.UUID(sqlcPost.UserID.Bytes),
			Title:       sqlcPost.Title,
			Content:     sqlcPost.Content.String,
			Slug:        sqlcPost.Slug,
			IsPublished: sqlcPost.IsPublished.Bool,
			PublishedAt: r.pgtypeTimestampToTimePtr(sqlcPost.PublishedAt),
			CreatedAt:   sqlcPost.CreatedAt.Time,
			UpdatedAt:   sqlcPost.UpdatedAt.Time,
		},
		AuthorEmail:     sqlcPost.AuthorEmail,
		AuthorFirstName: sqlcPost.AuthorFirstName,
		AuthorLastName:  sqlcPost.AuthorLastName,
	}
}

func (r *postRepository) mapSQLCGetPostBySlugRowToEntity(sqlcPost sqldb.GetPostBySlugRow) *entities.PostWithAuthor {
	return &entities.PostWithAuthor{
		Post: entities.Post{
			ID:          uuid.UUID(sqlcPost.ID.Bytes),
			UserID:      uuid.UUID(sqlcPost.UserID.Bytes),
			Title:       sqlcPost.Title,
			Content:     sqlcPost.Content.String,
			Slug:        sqlcPost.Slug,
			IsPublished: sqlcPost.IsPublished.Bool,
			PublishedAt: r.pgtypeTimestampToTimePtr(sqlcPost.PublishedAt),
			CreatedAt:   sqlcPost.CreatedAt.Time,
			UpdatedAt:   sqlcPost.UpdatedAt.Time,
		},
		AuthorEmail:     sqlcPost.AuthorEmail,
		AuthorFirstName: sqlcPost.AuthorFirstName,
		AuthorLastName:  sqlcPost.AuthorLastName,
	}
}

func (r *postRepository) mapSQLCListPostsRowToEntity(sqlcPost sqldb.ListPostsRow) *entities.PostWithAuthor {
	return &entities.PostWithAuthor{
		Post: entities.Post{
			ID:          uuid.UUID(sqlcPost.ID.Bytes),
			UserID:      uuid.UUID(sqlcPost.UserID.Bytes),
			Title:       sqlcPost.Title,
			Content:     sqlcPost.Content.String,
			Slug:        sqlcPost.Slug,
			IsPublished: sqlcPost.IsPublished.Bool,
			PublishedAt: r.pgtypeTimestampToTimePtr(sqlcPost.PublishedAt),
			CreatedAt:   sqlcPost.CreatedAt.Time,
			UpdatedAt:   sqlcPost.UpdatedAt.Time,
		},
		AuthorEmail:     sqlcPost.AuthorEmail,
		AuthorFirstName: sqlcPost.AuthorFirstName,
		AuthorLastName:  sqlcPost.AuthorLastName,
	}
}

func (r *postRepository) pgtypeTimestampToTimePtr(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
