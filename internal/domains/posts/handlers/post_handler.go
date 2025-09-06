package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/models"
	"github.com/zercle/gofiber-skeleton/internal/domains/posts/usecases"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
	"github.com/zercle/gofiber-skeleton/pkg/utils"
)

type PostHandler struct {
	postUsecase usecases.PostUsecase
}

func NewPostHandler(postUsecase usecases.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePostRequest true "Create post request"
// @Success 201 {object} models.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendFail(c, fiber.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return utils.SendValidationError(c, validationErrors)
	}

	result, err := h.postUsecase.CreatePost(c.Context(), userID, req)
	if err != nil {
		return utils.SendInternalError(c, "Failed to create post")
	}

	return utils.SendCreated(c, result)
}

// @Summary Get post by ID
// @Description Get post information by ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.PostResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	postID := c.Params("id")

	result, err := h.postUsecase.GetPost(c.Context(), postID)
	if err != nil {
		switch err {
		case types.ErrNotFound:
			return utils.SendNotFound(c, "Post not found")
		default:
			return utils.SendInternalError(c, "Failed to get post")
		}
	}

	return utils.SendSuccess(c, result)
}

// @Summary Update post
// @Description Update post information
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body models.UpdatePostRequest true "Update post request"
// @Success 200 {object} models.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	postID := c.Params("id")

	var req models.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendFail(c, fiber.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return utils.SendValidationError(c, validationErrors)
	}

	result, err := h.postUsecase.UpdatePost(c.Context(), postID, userID, req)
	if err != nil {
		switch err {
		case types.ErrNotFound:
			return utils.SendNotFound(c, "Post not found")
		case types.ErrForbidden:
			return utils.SendForbidden(c, "You can only update your own posts")
		default:
			return utils.SendInternalError(c, "Failed to update post")
		}
	}

	return utils.SendSuccess(c, result)
}

// @Summary Delete post
// @Description Delete a post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	postID := c.Params("id")

	err := h.postUsecase.DeletePost(c.Context(), postID, userID)
	if err != nil {
		switch err {
		case types.ErrNotFound:
			return utils.SendNotFound(c, "Post not found")
		case types.ErrForbidden:
			return utils.SendForbidden(c, "You can only delete your own posts")
		default:
			return utils.SendInternalError(c, "Failed to delete post")
		}
	}

	return utils.SendSuccess(c, map[string]interface{}{
		"message": "Post deleted successfully",
	})
}

// @Summary List posts
// @Description Get list of posts
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param published_only query bool false "Show only published posts" default(true)
// @Success 200 {object} models.PostsListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /posts [get]
func (h *PostHandler) ListPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	publishedOnly := c.Query("published_only", "true") == "true"

	result, err := h.postUsecase.ListPosts(c.Context(), page, pageSize, publishedOnly)
	if err != nil {
		return utils.SendInternalError(c, "Failed to list posts")
	}

	return utils.SendSuccess(c, result)
}

// @Summary List user posts
// @Description Get list of current user's posts
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param published_only query bool false "Show only published posts" default(false)
// @Success 200 {object} models.PostsListResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/me [get]
func (h *PostHandler) ListUserPosts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	publishedOnly := c.Query("published_only", "false") == "true"

	result, err := h.postUsecase.ListUserPosts(c.Context(), userID, page, pageSize, publishedOnly)
	if err != nil {
		return utils.SendInternalError(c, "Failed to list user posts")
	}

	return utils.SendSuccess(c, result)
}

// @Summary Publish post
// @Description Publish a post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} models.PostResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/publish [post]
func (h *PostHandler) PublishPost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	postID := c.Params("id")

	result, err := h.postUsecase.PublishPost(c.Context(), postID, userID)
	if err != nil {
		switch err {
		case types.ErrNotFound:
			return utils.SendNotFound(c, "Post not found")
		case types.ErrForbidden:
			return utils.SendForbidden(c, "You can only publish your own posts")
		default:
			return utils.SendInternalError(c, "Failed to publish post")
		}
	}

	return utils.SendSuccess(c, result)
}

// @Summary Unpublish post
// @Description Unpublish a post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} models.PostResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/unpublish [post]
func (h *PostHandler) UnpublishPost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	postID := c.Params("id")

	result, err := h.postUsecase.UnpublishPost(c.Context(), postID, userID)
	if err != nil {
		switch err {
		case types.ErrNotFound:
			return utils.SendNotFound(c, "Post not found")
		case types.ErrForbidden:
			return utils.SendForbidden(c, "You can only unpublish your own posts")
		default:
			return utils.SendInternalError(c, "Failed to unpublish post")
		}
	}

	return utils.SendSuccess(c, result)
}