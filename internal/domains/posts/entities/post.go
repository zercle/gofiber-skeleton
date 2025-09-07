package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
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

type PostWithAuthor struct {
	Post
	AuthorEmail     string `json:"author_email"`
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName  string `json:"author_last_name"`
}

func (p *PostWithAuthor) AuthorFullName() string {
	return p.AuthorFirstName + " " + p.AuthorLastName
}
