package models

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=10"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=10"`
}

type PublishPostRequest struct {
	IsPublished bool `json:"is_published" validate:"required"`
}