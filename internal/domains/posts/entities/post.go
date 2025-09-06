package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	AuthorID    string     `json:"author_id" db:"author_id"`
	IsPublished bool       `json:"is_published" db:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func NewPost(title, content, authorID string) (*Post, error) {
	now := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Post{
		ID:          id.String(),
		Title:       title,
		Content:     content,
		AuthorID:    authorID,
		IsPublished: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (p *Post) Publish() {
	if !p.IsPublished {
		now := time.Now()
		p.IsPublished = true
		p.PublishedAt = &now
		p.UpdatedAt = now
	}
}

func (p *Post) Unpublish() {
	if p.IsPublished {
		p.IsPublished = false
		p.PublishedAt = nil
		p.UpdatedAt = time.Now()
	}
}

func (p *Post) Update(title, content string) {
	p.Title = title
	p.Content = content
	p.UpdatedAt = time.Now()
}

func (p *Post) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.UpdatedAt = now
	if p.IsPublished {
		p.Unpublish()
	}
}

func (p *Post) IsDeleted() bool {
	return p.DeletedAt != nil
}