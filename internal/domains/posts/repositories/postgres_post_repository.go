package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/entities"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
)

type PostgresPostRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPostRepository(db *pgxpool.Pool) PostRepository {
	return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) Create(ctx context.Context, post *entities.Post) error {
	query := `
		INSERT INTO posts (id, title, content, author_id, is_published, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	_, err := r.db.Exec(ctx, query,
		post.ID, post.Title, post.Content, post.AuthorID,
		post.IsPublished, post.PublishedAt, post.CreatedAt, post.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	
	return nil
}

func (r *PostgresPostRepository) GetByID(ctx context.Context, id string) (*entities.Post, error) {
	query := `
		SELECT id, title, content, author_id, is_published, published_at, created_at, updated_at, deleted_at
		FROM posts 
		WHERE id = $1 AND deleted_at IS NULL`
	
	post := &entities.Post{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID,
		&post.IsPublished, &post.PublishedAt, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	
	return post, nil
}

func (r *PostgresPostRepository) Update(ctx context.Context, post *entities.Post) error {
	query := `
		UPDATE posts 
		SET title = $2, content = $3, is_published = $4, published_at = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL`
	
	result, err := r.db.Exec(ctx, query,
		post.ID, post.Title, post.Content, post.IsPublished, post.PublishedAt, post.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return types.ErrNotFound
	}
	
	return nil
}

func (r *PostgresPostRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE posts 
		SET deleted_at = NOW(), updated_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return types.ErrNotFound
	}
	
	return nil
}

func (r *PostgresPostRepository) List(ctx context.Context, limit, offset int, publishedOnly bool) ([]*entities.Post, error) {
	query := `
		SELECT id, title, content, author_id, is_published, published_at, created_at, updated_at, deleted_at
		FROM posts 
		WHERE deleted_at IS NULL`
	
	if publishedOnly {
		query += " AND is_published = true"
	}
	
	query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()
	
	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID,
			&post.IsPublished, &post.PublishedAt, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	
	return posts, nil
}

func (r *PostgresPostRepository) ListByAuthor(ctx context.Context, authorID string, limit, offset int, publishedOnly bool) ([]*entities.Post, error) {
	query := `
		SELECT id, title, content, author_id, is_published, published_at, created_at, updated_at, deleted_at
		FROM posts 
		WHERE author_id = $1 AND deleted_at IS NULL`
	
	if publishedOnly {
		query += " AND is_published = true"
	}
	
	query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by author: %w", err)
	}
	defer rows.Close()
	
	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID,
			&post.IsPublished, &post.PublishedAt, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	
	return posts, nil
}

func (r *PostgresPostRepository) Count(ctx context.Context, publishedOnly bool) (int64, error) {
	query := `SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL`
	
	if publishedOnly {
		query += " AND is_published = true"
	}
	
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}
	
	return count, nil
}

func (r *PostgresPostRepository) CountByAuthor(ctx context.Context, authorID string, publishedOnly bool) (int64, error) {
	query := `SELECT COUNT(*) FROM posts WHERE author_id = $1 AND deleted_at IS NULL`
	
	if publishedOnly {
		query += " AND is_published = true"
	}
	
	var count int64
	err := r.db.QueryRow(ctx, query, authorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts by author: %w", err)
	}
	
	return count, nil
}