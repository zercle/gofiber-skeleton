package models

import (
	"time"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/entities"
)

type PostResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	AuthorID    string     `json:"author_id"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PostsListResponse struct {
	Posts      []PostResponse `json:"posts"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

func NewPostResponse(post *entities.Post) PostResponse {
	return PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		AuthorID:    post.AuthorID,
		IsPublished: post.IsPublished,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}

func NewPostsListResponse(posts []*entities.Post, total int64, page, pageSize int) PostsListResponse {
	postResponses := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		postResponses = append(postResponses, NewPostResponse(post))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return PostsListResponse{
		Posts:      postResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}