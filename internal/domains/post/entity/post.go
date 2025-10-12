package entity

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zercle/gofiber-skeleton/pkg/uuid"
	uuidv7 "github.com/google/uuid"
)

var (
	ErrPostNotFound        = errors.New("post not found")
	ErrPostTitleRequired   = errors.New("post title is required")
	ErrPostContentRequired = errors.New("post content is required")
	ErrInvalidPostStatus   = errors.New("invalid post status")
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type DomainPost struct {
	ID        uuidv7.UUID  `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Status    PostStatus   `json:"status"`
	UserID    uuidv7.UUID  `json:"user_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (p *DomainPost) Validate() error {
	if p.Title == "" {
		return ErrPostTitleRequired
	}
	if p.Content == "" {
		return ErrPostContentRequired
	}
	if p.Status != PostStatusDraft && p.Status != PostStatusPublished && p.Status != PostStatusArchived {
		return ErrInvalidPostStatus
	}
	return nil
}

func (p *DomainPost) IsPublished() bool {
	return p.Status == PostStatusPublished
}

func (p *DomainPost) IsDraft() bool {
	return p.Status == PostStatusDraft
}

func (p *DomainPost) CanBeAccessedBy(userID uuidv7.UUID) bool {
	return p.UserID == userID
}

type PostWithAuthor struct {
	ID        uuidv7.UUID  `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Status    PostStatus   `json:"status"`
	UserID    uuidv7.UUID  `json:"user_id"`
	FullName  string       `json:"full_name"`
	Email     string       `json:"email"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type PostStats struct {
	TotalPosts     int        `json:"total_posts"`
	PublishedPosts int        `json:"published_posts"`
	DraftPosts     int        `json:"draft_posts"`
	LastPostDate   *time.Time `json:"last_post_date"`
}

type PostStatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type CreatePostRequest struct {
	Title   string     `json:"title" validate:"required,min=1,max=255"`
	Content string     `json:"content" validate:"required,min=1"`
	Status  PostStatus `json:"status" validate:"oneof=draft published archived"`
}

type UpdatePostRequest struct {
	Title   *string     `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content *string     `json:"content,omitempty" validate:"omitempty,min=1"`
	Status  *PostStatus `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

type PostFilter struct {
	Status *PostStatus   `json:"status,omitempty"`
	UserID *uuidv7.UUID  `json:"user_id,omitempty"`
	Search *string       `json:"search,omitempty"`
	Limit  int           `json:"limit" validate:"min=1,max=100"`
	Offset int           `json:"offset" validate:"min=0"`
}

func NewPost(title, content string, userID uuidv7.UUID) *DomainPost {
	return &DomainPost{
		ID:        uuid.NewV7(),
		Title:     title,
		Content:   content,
		Status:    PostStatusDraft,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Convert between Post (sqlc generated) and DomainPost
func PostToDomainPost(sqlcPost Post) *DomainPost {
	return &DomainPost{
		ID:        uuidv7.UUID(sqlcPost.ID.Bytes),
		Title:     sqlcPost.Title,
		Content:   sqlcPost.Content,
		Status:    PostStatus(sqlcPost.Status),
		UserID:    uuidv7.UUID(sqlcPost.UserID.Bytes),
		CreatedAt: sqlcPost.CreatedAt.Time,
		UpdatedAt: sqlcPost.UpdatedAt.Time,
	}
}

func DomainPostToPost(p *DomainPost) Post {
	var postUUID [16]byte
	var userUUID [16]byte
	copy(postUUID[:], p.ID[:])
	copy(userUUID[:], p.UserID[:])
	return Post{
		ID:        pgtype.UUID{Bytes: postUUID, Valid: true},
		Title:     p.Title,
		Content:   p.Content,
		Status:    string(p.Status),
		UserID:    pgtype.UUID{Bytes: userUUID, Valid: true},
		CreatedAt: pgtype.Timestamptz{Time: p.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: p.UpdatedAt, Valid: true},
	}
}

func DefaultPostFilter() *PostFilter {
	return &PostFilter{
		Limit:  10,
		Offset: 0,
	}
}
