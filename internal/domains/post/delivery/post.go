package delivery

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/usecase"
	"github.com/zercle/gofiber-skeleton/pkg/response"
	"github.com/zercle/gofiber-skeleton/pkg/utils"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
}

func NewPostHandler(injector do.Injector) (*PostHandler, error) {
	postUsecase := do.MustInvoke[usecase.PostUsecase](injector)
	return &PostHandler{
		postUsecase: postUsecase,
	}, nil
}

// RegisterRoutes registers post routes
func (h *PostHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Public post routes
	api.Get("/posts", h.GetPublishedPosts)
	api.Get("/posts/search", h.SearchPosts)
	api.Get("/posts/:id", h.GetPost)

	// Protected post routes (require authentication)
	protected := api.Group("/posts")
	protected.Use(utils.AuthMiddleware())
	protected.Post("/", h.CreatePost)
	protected.Put("/:id", h.UpdatePost)
	protected.Delete("/:id", h.DeletePost)
	protected.Post("/:id/publish", h.PublishPost)
	protected.Post("/:id/archive", h.ArchivePost)
	protected.Post("/:id/unpublish", h.UnpublishPost)

	// User-specific routes
	userPosts := api.Group("/users/:userId/posts")
	userPosts.Get("/", h.GetUserPosts)
	userPosts.Get("/stats", h.GetUserPostStats)
}

// CreatePost creates a new post
// @Summary Create a new post
// @Description Create a new post for the authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param post body entity.CreatePostRequest true "Post data"
// @Security ApiKeyAuth
// @Success 201 {object} response.Response{data=entity.Post}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var req entity.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Invalid user ID", err.Error())
	}

	post, err := h.postUsecase.CreatePost(c.Context(), userID, &req)
	if err != nil {
		if err == entity.ErrPostTitleRequired || err == entity.ErrPostContentRequired || err == entity.ErrInvalidPostStatus {
			return response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to create post", err.Error())
	}

	return response.Success(c, http.StatusCreated, "Post created successfully", post)
}

// GetPost retrieves a post by ID
// @Summary Get a post by ID
// @Description Retrieve a post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response{data=entity.Post}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts/{id} [get]
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid post ID", err.Error())
	}

	post, err := h.postUsecase.GetPost(c.Context(), postID)
	if err != nil {
		if err == entity.ErrPostNotFound {
			return response.Error(c, http.StatusNotFound, "Post not found", err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to get post", err.Error())
	}

	return response.Success(c, http.StatusOK, "Post retrieved successfully", post)
}

// UpdatePost updates an existing post
// @Summary Update a post
// @Description Update an existing post (only by owner)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body entity.UpdatePostRequest true "Post update data"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=entity.Post}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid post ID", err.Error())
	}

	var req entity.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Get user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Invalid user ID", err.Error())
	}

	post, err := h.postUsecase.UpdatePost(c.Context(), postID, userID, &req)
	if err != nil {
		if err == entity.ErrPostNotFound {
			return response.Error(c, http.StatusNotFound, "Post not found", err.Error())
		}
		if err == entity.ErrPostTitleRequired || err == entity.ErrPostContentRequired || err == entity.ErrInvalidPostStatus {
			return response.Error(c, http.StatusBadRequest, err.Error(), err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to update post", err.Error())
	}

	return response.Success(c, http.StatusOK, "Post updated successfully", post)
}

// DeletePost deletes a post
// @Summary Delete a post
// @Description Delete a post (only by owner)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid post ID", err.Error())
	}

	// Get user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Invalid user ID", err.Error())
	}

	err = h.postUsecase.DeletePost(c.Context(), postID, userID)
	if err != nil {
		if err == entity.ErrPostNotFound {
			return response.Error(c, http.StatusNotFound, "Post not found", err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to delete post", err.Error())
	}

	return response.Success(c, http.StatusOK, "Post deleted successfully", nil)
}

// GetPublishedPosts retrieves published posts with pagination
// @Summary Get published posts
// @Description Retrieve a list of published posts with pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of posts" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} response.Response{data=[]entity.PostWithAuthor}
// @Failure 500 {object} response.Response
// @Router /api/v1/posts [get]
func (h *PostHandler) GetPublishedPosts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	posts, err := h.postUsecase.GetPublishedPosts(c.Context(), limit, offset)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get posts", err.Error())
	}

	return response.Success(c, http.StatusOK, "Posts retrieved successfully", posts)
}

// GetUserPosts retrieves posts for a specific user
// @Summary Get user posts
// @Description Retrieve posts for a specific user with optional status filter
// @Tags posts
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param status query string false "Post status filter (draft, published, archived)"
// @Success 200 {object} response.Response{data=[]entity.Post}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{userId}/posts [get]
func (h *PostHandler) GetUserPosts(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
	}

	var status *entity.PostStatus
	if statusStr := c.Query("status"); statusStr != "" {
		postStatus := entity.PostStatus(statusStr)
		if postStatus != entity.PostStatusDraft && postStatus != entity.PostStatusPublished && postStatus != entity.PostStatusArchived {
			return response.Error(c, http.StatusBadRequest, "Invalid status filter", "Invalid status filter")
		}
		status = &postStatus
	}

	posts, err := h.postUsecase.GetUserPosts(c.Context(), userID, status)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get user posts", err.Error())
	}

	return response.Success(c, http.StatusOK, "User posts retrieved successfully", posts)
}

// GetUserPostStats retrieves post statistics for a user
// @Summary Get user post statistics
// @Description Retrieve post statistics for a specific user
// @Tags posts
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} response.Response{data=entity.PostStats}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{userId}/posts/stats [get]
func (h *PostHandler) GetUserPostStats(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
	}

	stats, err := h.postUsecase.GetUserPostStats(c.Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get user post stats", err.Error())
	}

	return response.Success(c, http.StatusOK, "User post statistics retrieved successfully", stats)
}

// SearchPosts searches published posts
// @Summary Search posts
// @Description Search published posts by title or content
// @Tags posts
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit number of posts" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} response.Response{data=[]entity.PostWithAuthor}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts/search [get]
func (h *PostHandler) SearchPosts(c *fiber.Ctx) error {
	searchQuery := c.Query("q")
	if searchQuery == "" {
		return response.Error(c, http.StatusBadRequest, "Search query is required", "Search query is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	posts, err := h.postUsecase.SearchPosts(c.Context(), searchQuery, limit, offset)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to search posts", err.Error())
	}

	return response.Success(c, http.StatusOK, "Posts searched successfully", posts)
}

// PublishPost publishes a post
// @Summary Publish a post
// @Description Publish a post (only by owner)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts/{id}/publish [post]
func (h *PostHandler) PublishPost(c *fiber.Ctx) error {
	return h.updatePostStatus(c, "publish")
}

// ArchivePost archives a post
// @Summary Archive a post
// @Description Archive a post (only by owner)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts/{id}/archive [post]
func (h *PostHandler) ArchivePost(c *fiber.Ctx) error {
	return h.updatePostStatus(c, "archive")
}

// UnpublishPost unpublishes a post
// @Summary Unpublish a post
// @Description Unpublish a post (sets to draft, only by owner)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/posts/{id}/unpublish [post]
func (h *PostHandler) UnpublishPost(c *fiber.Ctx) error {
	return h.updatePostStatus(c, "unpublish")
}

// updatePostStatus is a helper method for post status updates
func (h *PostHandler) updatePostStatus(c *fiber.Ctx, action string) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid post ID", err.Error())
	}

	// Get user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Invalid user ID", err.Error())
	}

	var errFunc func(context.Context, uuid.UUID, uuid.UUID) error
	var successMsg string

	switch action {
	case "publish":
		errFunc = h.postUsecase.PublishPost
		successMsg = "Post published successfully"
	case "archive":
		errFunc = h.postUsecase.ArchivePost
		successMsg = "Post archived successfully"
	case "unpublish":
		errFunc = h.postUsecase.UnpublishPost
		successMsg = "Post unpublished successfully"
	default:
		return response.Error(c, http.StatusBadRequest, "Invalid action", "Invalid action")
	}

	err = errFunc(c.Context(), postID, userID)
	if err != nil {
		if err == entity.ErrPostNotFound {
			return response.Error(c, http.StatusNotFound, "Post not found", err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to update post status", err.Error())
	}

	return response.Success(c, http.StatusOK, successMsg, nil)
}
