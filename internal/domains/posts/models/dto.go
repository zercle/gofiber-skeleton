package models

import (
	"time"

	"github.com/google/uuid"
)

// CreatePostRequest represents the request payload for creating a post.
type CreatePostRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Content     string `json:"content" validate:"required,min=10"`
	IsPublished bool   `json:"is_published"`
}

// UpdatePostRequest represents the request payload for updating a post.
type UpdatePostRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Content     *string `json:"content,omitempty" validate:"omitempty,min=10"`
	IsPublished *bool   `json:"is_published,omitempty"`
}

// PostResponse represents the response payload for a post.
type PostResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Slug        string     `json:"slug"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// PostWithAuthorResponse represents the response payload for a post with author info.
type PostWithAuthorResponse struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Slug            string     `json:"slug"`
	IsPublished     bool       `json:"is_published"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	AuthorEmail     string     `json:"author_email"`
	AuthorFirstName string     `json:"author_first_name"`
	AuthorLastName  string     `json:"author_last_name"`
	AuthorFullName  string     `json:"author_full_name"`
}

// PostListResponse represents the response payload for listing posts.
type PostListResponse struct {
	Posts  []PostWithAuthorResponse `json:"posts"`
	Total  int                      `json:"total"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
}
