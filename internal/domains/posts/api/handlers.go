package api

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/internal/domains/posts/models"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/usecases"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/validation"
	"github.com/zercle/gofiber-skeleton/internal/shared/jsend"
)

type PostHandler struct {
	postUseCase usecases.PostUseCase
	validator   *validator.Validate
}

func NewPostHandler(postUseCase usecases.PostUseCase) *PostHandler {
	return &PostHandler{
		postUseCase: postUseCase,
		validator:   validation.NewValidator(),
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePostRequest true "Post data"
// @Success 201 {object} jsend.JSendResponse{data=models.PostResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	var req models.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		var errors []jsend.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, jsend.ValidationError{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		return jsend.SendValidationError(c, errors)
	}

	post, err := h.postUseCase.Create(c.Context(), userID, &req)
	if err != nil {
		return jsend.SendInternalServerError(c, err.Error())
	}

	return jsend.SendSuccess(c.Status(fiber.StatusCreated), post)
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Retrieve a specific post by its ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} jsend.JSendResponse{data=models.PostWithAuthorResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 404 {object} jsend.JSendResponse
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid post ID",
		})
	}

	post, err := h.postUseCase.GetByID(c.Context(), id)
	if err != nil {
		return jsend.SendNotFound(c, map[string]string{
			"message": err.Error(),
		})
	}

	return jsend.SendSuccess(c, post)
}

// GetPostBySlug godoc
// @Summary Get a post by slug
// @Description Retrieve a specific post by its slug
// @Tags posts
// @Produce json
// @Param slug path string true "Post Slug"
// @Success 200 {object} jsend.JSendResponse{data=models.PostWithAuthorResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 404 {object} jsend.JSendResponse
// @Router /posts/slug/{slug} [get]
func (h *PostHandler) GetPostBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Slug is required",
		})
	}

	post, err := h.postUseCase.GetBySlug(c.Context(), slug)
	if err != nil {
		return jsend.SendNotFound(c, map[string]string{
			"message": err.Error(),
		})
	}

	return jsend.SendSuccess(c, post)
}

// ListPosts godoc
// @Summary List posts
// @Description Retrieve a paginated list of posts
// @Tags posts
// @Produce json
// @Param limit query int false "Number of posts to return" default(20)
// @Param offset query int false "Number of posts to skip" default(0)
// @Param published query bool false "Filter by published status"
// @Success 200 {object} jsend.JSendResponse{data=models.PostListResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Router /posts [get]
func (h *PostHandler) ListPosts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	var isPublished *bool
	if publishedStr := c.Query("published"); publishedStr != "" {
		if published, err := strconv.ParseBool(publishedStr); err == nil {
			isPublished = &published
		}
	}

	if limit > 100 {
		limit = 100 // Prevent excessive requests
	}

	posts, err := h.postUseCase.List(c.Context(), limit, offset, isPublished)
	if err != nil {
		return jsend.SendInternalServerError(c, err.Error())
	}

	return jsend.SendSuccess(c, posts)
}

// ListUserPosts godoc
// @Summary List user's posts
// @Description Retrieve a paginated list of the authenticated user's posts
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of posts to return" default(20)
// @Param offset query int false "Number of posts to skip" default(0)
// @Success 200 {object} jsend.JSendResponse{data=models.PostListResponse}
// @Failure 401 {object} jsend.JSendResponse
// @Router /posts/my-posts [get]
func (h *PostHandler) ListUserPosts(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100 // Prevent excessive requests
	}

	posts, err := h.postUseCase.ListByUser(c.Context(), userID, limit, offset)
	if err != nil {
		return jsend.SendInternalServerError(c, err.Error())
	}

	return jsend.SendSuccess(c, posts)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body models.UpdatePostRequest true "Update post data"
// @Success 200 {object} jsend.JSendResponse{data=models.PostResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Failure 403 {object} jsend.JSendResponse
// @Failure 404 {object} jsend.JSendResponse
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid post ID",
		})
	}

	var req models.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		var errors []jsend.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, jsend.ValidationError{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		return jsend.SendValidationError(c, errors)
	}

	post, err := h.postUseCase.Update(c.Context(), id, userID, &req)
	if err != nil {
		if err.Error() == "post not found" {
			return jsend.SendNotFound(c, map[string]string{
				"message": err.Error(),
			})
		}
		if err.Error() == "unauthorized: you can only update your own posts" {
			return jsend.SendForbidden(c, map[string]string{
				"message": err.Error(),
			})
		}
		return jsend.SendInternalServerError(c, err.Error())
	}

	return jsend.SendSuccess(c, post)
}

// PublishPost godoc
// @Summary Publish a post
// @Description Publish an existing post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} jsend.JSendResponse{data=models.PostResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Failure 403 {object} jsend.JSendResponse
// @Failure 404 {object} jsend.JSendResponse
// @Router /posts/{id}/publish [patch]
func (h *PostHandler) PublishPost(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid post ID",
		})
	}

	post, err := h.postUseCase.Publish(c.Context(), id, userID)
	if err != nil {
		if err.Error() == "post not found" {
			return jsend.SendNotFound(c, map[string]string{
				"message": err.Error(),
			})
		}
		if err.Error() == "unauthorized: you can only publish your own posts" {
			return jsend.SendForbidden(c, map[string]string{
				"message": err.Error(),
			})
		}
		return jsend.SendInternalServerError(c, err.Error())
	}

	return jsend.SendSuccess(c, post)
}

// UnpublishPost godoc
// @Summary Unpublish a post
// @Description Unpublish an existing post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} jsend.JSendResponse{data=models.PostResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Failure 403 {object} jsend.JSendResponse
// @Failure 404 {object} jsend.JSendResponse
// @Router /posts/{id}/unpublish [patch]
func (h *PostHandler) UnpublishPost(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid post ID",
		})
	}

	post, err := h.postUseCase.Unpublish(c.Context(), id, userID)
	if err != nil {
		if err.Error() == "post not found" {
			return jsend.SendNotFound(c, map[string]string{
				"message": err.Error(),
			})
		}
		if err.Error() == "unauthorized: you can only unpublish your own posts" {
			return jsend.SendForbidden(c, map[string]string{
				"message": err.Error(),
			})
		}
		return jsend.SendInternalServerError(c, err.Error())
	}

	return jsend.SendSuccess(c, post)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete an existing post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 204
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Failure 403 {object} jsend.JSendResponse
// @Failure 404 {object} jsend.JSendResponse
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid post ID",
		})
	}

	err = h.postUseCase.Delete(c.Context(), id, userID)
	if err != nil {
		if err.Error() == "post not found" {
			return jsend.SendNotFound(c, map[string]string{
				"message": err.Error(),
			})
		}
		if err.Error() == "unauthorized: you can only delete your own posts" {
			return jsend.SendForbidden(c, map[string]string{
				"message": err.Error(),
			})
		}
		return jsend.SendInternalServerError(c, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}