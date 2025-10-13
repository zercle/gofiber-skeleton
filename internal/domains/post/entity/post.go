package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=10"`
	Status  string `json:"status" validate:"required,oneof=draft published"`
}

type UpdatePostRequest struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Content *string `json:"content,omitempty" validate:"omitempty,min=10"`
	Status  *string `json:"status,omitempty" validate:"omitempty,oneof=draft published"`
}

type PostResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  uuid.UUID `json:"author_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	PostStatusDraft     = "draft"
	PostStatusPublished = "published"
)

func NewPost(title, content, status string, authorID uuid.UUID) *Post {
	now := time.Now()
	return &Post{
		ID:        uuid.New(),
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *Post) ToResponse() *PostResponse {
	return &PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *Post) Update(title, content, status string) {
	p.Title = title
	p.Content = content
	p.Status = status
	p.UpdatedAt = time.Now()
}

func IsValidStatus(status string) bool {
	return status == PostStatusDraft || status == PostStatusPublished
}