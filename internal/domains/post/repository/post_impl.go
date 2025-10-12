package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/pkg/database"
)

type sqlcPostRepository struct {
	db      *database.Database
	queries *entity.Queries
}

// NewPostRepository creates a new post repository implementation
func NewPostRepository(injector do.Injector) (PostRepository, error) {
	db := do.MustInvoke[*database.Database](injector)
	return &sqlcPostRepository{
		db:      db,
		queries: entity.New(db.GetPool()),
	}, nil
}

// Create creates a new post
func (r *sqlcPostRepository) Create(ctx context.Context, post *entity.DomainPost) error {
	var userUUID [16]byte
	copy(userUUID[:], post.UserID[:])
	userPgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	var postUUID [16]byte
	copy(postUUID[:], post.ID[:])
	postPgUUID := pgtype.UUID{Bytes: postUUID, Valid: true}

	params := entity.CreatePostParams{
		ID:        postPgUUID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    string(post.Status),
		UserID:    userPgUUID,
		CreatedAt: pgtype.Timestamptz{Time: post.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}

	_, err := r.queries.CreatePost(ctx, params)
	return err
}

// GetByID retrieves a post by ID
func (r *sqlcPostRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.DomainPost, error) {
	var postUUID [16]byte
	copy(postUUID[:], id[:])
	pgUUID := pgtype.UUID{Bytes: postUUID, Valid: true}

	post, err := r.queries.GetPostByID(ctx, pgUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrPostNotFound
		}
		return nil, err
	}

	return &entity.DomainPost{
		ID:        uuid.UUID(post.ID.Bytes),
		Title:     post.Title,
		Content:   post.Content,
		Status:    entity.PostStatus(post.Status),
		UserID:    uuid.UUID(post.UserID.Bytes),
		CreatedAt: post.CreatedAt.Time,
		UpdatedAt: post.UpdatedAt.Time,
	}, nil
}

// Update updates a post
func (r *sqlcPostRepository) Update(ctx context.Context, post *entity.DomainPost) error {
	var postUUID [16]byte
	copy(postUUID[:], post.ID[:])
	pgUUID := pgtype.UUID{Bytes: postUUID, Valid: true}

	params := entity.UpdatePostParams{
		ID:        pgUUID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    string(post.Status),
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}

	_, err := r.queries.UpdatePost(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.ErrPostNotFound
		}
		return err
	}

	return nil
}

// Delete deletes a post
func (r *sqlcPostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var postUUID [16]byte
	copy(postUUID[:], id[:])
	pgUUID := pgtype.UUID{Bytes: postUUID, Valid: true}

	err := r.queries.DeletePost(ctx, pgUUID)
	if err != nil {
		return err
	}

	return nil
}

// GetByUserID retrieves posts by user ID
func (r *sqlcPostRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.DomainPost, error) {
	var userUUID [16]byte
	copy(userUUID[:], userID[:])
	userPgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	posts, err := r.queries.GetPostsByUserID(ctx, userPgUUID)
	if err != nil {
		return nil, err
	}

	var result []*entity.DomainPost
	for _, post := range posts {
		result = append(result, &entity.DomainPost{
			ID:        uuid.UUID(post.ID.Bytes),
			Title:     post.Title,
			Content:   post.Content,
			Status:    entity.PostStatus(post.Status),
			UserID:    uuid.UUID(post.UserID.Bytes),
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	return result, nil
}

// GetPublishedPostsByUserID retrieves published posts by user ID
func (r *sqlcPostRepository) GetPublishedPostsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.DomainPost, error) {
	var userUUID [16]byte
	copy(userUUID[:], userID[:])
	userPgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	posts, err := r.queries.GetPublishedPostsByUserID(ctx, userPgUUID)
	if err != nil {
		return nil, err
	}

	var result []*entity.DomainPost
	for _, post := range posts {
		result = append(result, &entity.DomainPost{
			ID:        uuid.UUID(post.ID.Bytes),
			Title:     post.Title,
			Content:   post.Content,
			Status:    entity.PostStatus(post.Status),
			UserID:    uuid.UUID(post.UserID.Bytes),
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	return result, nil
}

// GetDraftPostsByUserID retrieves draft posts by user ID
func (r *sqlcPostRepository) GetDraftPostsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.DomainPost, error) {
	var userUUID [16]byte
	copy(userUUID[:], userID[:])
	userPgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	posts, err := r.queries.GetDraftPostsByUserID(ctx, userPgUUID)
	if err != nil {
		return nil, err
	}

	var result []*entity.DomainPost
	for _, post := range posts {
		result = append(result, &entity.DomainPost{
			ID:        uuid.UUID(post.ID.Bytes),
			Title:     post.Title,
			Content:   post.Content,
			Status:    entity.PostStatus(post.Status),
			UserID:    uuid.UUID(post.UserID.Bytes),
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	return result, nil
}

// GetAllPosts retrieves all posts with user info
func (r *sqlcPostRepository) GetAllPosts(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error) {
	params := entity.GetAllPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	posts, err := r.queries.GetAllPosts(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*entity.PostWithAuthor
	for _, post := range posts {
		result = append(result, &entity.PostWithAuthor{
			ID:        uuid.UUID(post.ID.Bytes),
			Title:     post.Title,
			Content:   post.Content,
			Status:    entity.PostStatus(post.Status),
			UserID:    uuid.UUID(post.UserID.Bytes),
			FullName:  post.FullName,
			Email:     post.Email,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	return result, nil
}

// GetPostsWithAuthor retrieves published posts with author info
func (r *sqlcPostRepository) GetPostsWithAuthor(ctx context.Context, limit, offset int) ([]*entity.PostWithAuthor, error) {
	params := entity.GetPostsWithAuthorParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	posts, err := r.queries.GetPostsWithAuthor(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*entity.PostWithAuthor
	for _, post := range posts {
		result = append(result, &entity.PostWithAuthor{
			ID:        uuid.UUID(post.ID.Bytes),
			Title:     post.Title,
			Content:   post.Content,
			Status:    entity.PostStatus(post.Status),
			UserID:    uuid.UUID(post.UserID.Bytes),
			FullName:  post.FullName,
			Email:     post.Email,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	return result, nil
}

// SearchPosts searches posts with filters
func (r *sqlcPostRepository) SearchPosts(ctx context.Context, search string, status entity.PostStatus, limit, offset int) ([]*entity.PostWithAuthor, error) {
	searchPattern := "%" + search + "%"
	statusStr := string(status)

	params := entity.SearchPostsParams{
		Title:  searchPattern,
		Status: statusStr,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	posts, err := r.queries.SearchPosts(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*entity.PostWithAuthor
	for _, post := range posts {
		result = append(result, &entity.PostWithAuthor{
			ID:        uuid.UUID(post.ID.Bytes),
			Title:     post.Title,
			Content:   post.Content,
			Status:    entity.PostStatus(post.Status),
			UserID:    uuid.UUID(post.UserID.Bytes),
			FullName:  post.FullName,
			Email:     post.Email,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	return result, nil
}

// GetUserPostStats retrieves post statistics for a user
func (r *sqlcPostRepository) GetUserPostStats(ctx context.Context, userID uuid.UUID) (*entity.PostStats, error) {
	var userUUID [16]byte
	copy(userUUID[:], userID[:])
	userPgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	stats, err := r.queries.GetUserPostStats(ctx, userPgUUID)
	if err != nil {
		return nil, err
	}

	var lastPostDate *time.Time
	if stats.LastPostDate != nil {
		if t, ok := stats.LastPostDate.(time.Time); ok {
			lastPostDate = &t
		}
	}

	return &entity.PostStats{
		TotalPosts:     int(stats.TotalPosts),
		PublishedPosts: int(stats.PublishedPosts),
		DraftPosts:     int(stats.DraftPosts),
		LastPostDate:   lastPostDate,
	}, nil
}

// CountPostsByStatus returns count of posts grouped by status
func (r *sqlcPostRepository) CountPostsByStatus(ctx context.Context) ([]*entity.PostStatusCount, error) {
	statusCounts, err := r.queries.CountPostsByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Since SQLC generated this as a single row query, return it as a single item
	result := []*entity.PostStatusCount{
		{
			Status: statusCounts.Status,
			Count:  int(statusCounts.Count),
		},
	}

	return result, nil
}

// Exists checks if a post exists
func (r *sqlcPostRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		if err == entity.ErrPostNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsOwner checks if user is the owner of the post
func (r *sqlcPostRepository) IsOwner(ctx context.Context, postID, userID uuid.UUID) (bool, error) {
	post, err := r.GetByID(ctx, postID)
	if err != nil {
		if err == entity.ErrPostNotFound {
			return false, nil
		}
		return false, err
	}
	return post.UserID == userID, nil
}
