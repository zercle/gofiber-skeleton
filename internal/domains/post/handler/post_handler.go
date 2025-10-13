package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/usecase"
	"github.com/zercle/gofiber-skeleton/internal/shared/middleware"
	"github.com/zercle/gofiber-skeleton/internal/shared/response"
	"github.com/zercle/gofiber-skeleton/internal/shared/validator"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
}

func NewPostHandler(postUsecase usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with the provided details
// @Tags posts
// @Accept json
// @Produce json
// @Param request body entity.CreatePostRequest true "Post data"
// @Security BearerAuth
// @Success 201 {object} response.Response{data=entity.PostResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid user context", err.Error())
	}

	var req entity.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	post, err := h.postUsecase.Create(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrInvalidStatus:
			return response.BadRequest(c, "Invalid post status", err.Error())
		default:
			return response.InternalServerError(c, "Failed to create post", err.Error())
		}
	}

	return response.Created(c, "Post created successfully", post.ToResponse())
}

// GetPost godoc
// @Summary Get post by ID
// @Description Get a post by its ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.PostResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/posts/{id} [get]
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid post ID", err.Error())
	}

	post, err := h.postUsecase.GetByID(c.Context(), postID)
	if err != nil {
		if err == usecase.ErrPostNotFound {
			return response.NotFound(c, "Post not found", err.Error())
		}
		return response.InternalServerError(c, "Failed to get post", err.Error())
	}

	return response.OK(c, "Post retrieved successfully", post.ToResponse())
}

// ListPosts godoc
// @Summary List posts
// @Description Get a paginated list of all posts
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]entity.PostResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/posts [get]
func (h *PostHandler) ListPosts(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	posts, err := h.postUsecase.List(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to list posts", err.Error())
	}

	postResponses := make([]*entity.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	return response.OK(c, "Posts retrieved successfully", postResponses)
}

// ListMyPosts godoc
// @Summary List current user's posts
// @Description Get a paginated list of posts created by the current user
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]entity.PostResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/posts/my [get]
func (h *PostHandler) ListMyPosts(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid user context", err.Error())
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	posts, err := h.postUsecase.ListByAuthor(c.Context(), userID, limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to list user posts", err.Error())
	}

	postResponses := make([]*entity.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	return response.OK(c, "Your posts retrieved successfully", postResponses)
}

// UpdatePost godoc
// @Summary Update post
// @Description Update a post by its ID (only author can update)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param request body entity.UpdatePostRequest true "Post update data"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.PostResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid user context", err.Error())
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid post ID", err.Error())
	}

	var req entity.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	post, err := h.postUsecase.Update(c.Context(), userID, postID, &req)
	if err != nil {
		switch err {
		case usecase.ErrPostNotFound:
			return response.NotFound(c, "Post not found", err.Error())
		case usecase.ErrUnauthorized:
			return response.Forbidden(c, "Unauthorized to update this post", err.Error())
		case usecase.ErrInvalidStatus:
			return response.BadRequest(c, "Invalid post status", err.Error())
		default:
			return response.InternalServerError(c, "Failed to update post", err.Error())
		}
	}

	return response.OK(c, "Post updated successfully", post.ToResponse())
}

// DeletePost godoc
// @Summary Delete post
// @Description Delete a post by its ID (only author can delete)
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Security BearerAuth
// @Success 204 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid user context", err.Error())
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid post ID", err.Error())
	}

	if err := h.postUsecase.Delete(c.Context(), userID, postID); err != nil {
		switch err {
		case usecase.ErrPostNotFound:
			return response.NotFound(c, "Post not found", err.Error())
		case usecase.ErrUnauthorized:
			return response.Forbidden(c, "Unauthorized to delete this post", err.Error())
		default:
			return response.InternalServerError(c, "Failed to delete post", err.Error())
		}
	}

	return response.NoContent(c, "Post deleted successfully")
}