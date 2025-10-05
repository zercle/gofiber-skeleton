package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/post/usecase"
	"github.com/zercle/gofiber-skeleton/internal/response"
	"github.com/zercle/gofiber-skeleton/internal/user/middleware"
	"github.com/zercle/gofiber-skeleton/pkg/validator"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
}

func NewPostHandler(postUsecase usecase.PostUsecase) *PostHandler {
	return &PostHandler{postUsecase: postUsecase}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post in a thread (requires authentication)
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createPostRequest true "Post details"
// @Success 201 {object} response.JSendResponse "Post created successfully"
// @Failure 400 {object} response.JSendResponse "Invalid request"
// @Failure 401 {object} response.JSendResponse "Unauthorized"
// @Failure 500 {object} response.JSendResponse "Internal server error"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	req := new(createPostRequest)
	if err := c.BodyParser(req); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse create post request")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid request body"})
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().Error().Msg("user_id not found or invalid in context")
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", 2001)
	}

	threadID, err := uuid.Parse(req.ThreadID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse thread ID from request")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid thread ID"})
	}

	post, err := h.postUsecase.Create(c.Context(), userUUID, threadID, req.Content)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to create post")
		return response.Error(c, http.StatusInternalServerError, "Failed to create post", 2002)
	}

	return response.Success(c, http.StatusCreated, fiber.Map{
		"id":         post.ID,
		"thread_id":  post.ThreadID,
		"user_id":    post.UserID,
		"content":    post.Content,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	})
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Retrieve a single post by its ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID (UUID)"
// @Success 200 {object} response.JSendResponse "Post retrieved successfully"
// @Failure 400 {object} response.JSendResponse "Invalid post ID"
// @Failure 404 {object} response.JSendResponse "Post not found"
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse post ID from params")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid post ID"})
	}

	post, err := h.postUsecase.Get(c.Context(), postID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to get post")
		return response.Fail(c, http.StatusNotFound, fiber.Map{"error": "Post not found"})
	}

	return response.Success(c, http.StatusOK, fiber.Map{
		"id":         post.ID,
		"thread_id":  post.ThreadID,
		"user_id":    post.UserID,
		"content":    post.Content,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	})
}

// ListPostsByUser godoc
// @Summary List posts by user
// @Description Retrieve all posts created by a specific user
// @Tags posts
// @Produce json
// @Param user_id path string true "User ID (UUID)"
// @Success 200 {object} response.JSendResponse "Posts retrieved successfully"
// @Failure 400 {object} response.JSendResponse "Invalid user ID"
// @Failure 500 {object} response.JSendResponse "Internal server error"
// @Router /users/{user_id}/posts [get]
func (h *PostHandler) ListPostsByUser(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userUUID, err := uuid.Parse(userIDParam)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse user ID from params")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid user ID"})
	}

	posts, err := h.postUsecase.ListByUser(c.Context(), userUUID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to list posts by user")
		return response.Error(c, http.StatusInternalServerError, "Failed to retrieve posts", 2003)
	}

	return response.Success(c, http.StatusOK, fiber.Map{"posts": posts})
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update an existing post (requires authentication and ownership)
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID (UUID)"
// @Param request body updatePostRequest true "Updated post content"
// @Success 200 {object} response.JSendResponse "Post updated successfully"
// @Failure 400 {object} response.JSendResponse "Invalid request"
// @Failure 401 {object} response.JSendResponse "Unauthorized"
// @Failure 500 {object} response.JSendResponse "Internal server error"
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse post ID from params")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid post ID"})
	}

	req := new(updatePostRequest)
	if parseErr := c.BodyParser(req); parseErr != nil {
		logger.GetLogger().Error().Err(parseErr).Msg("Failed to parse update post request")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid request body"})
	}

	if validateErr := validator.ValidateRequest(c, req); validateErr != nil {
		return validateErr
	}

	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().Error().Msg("user_id not found or invalid in context")
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", 2001)
	}

	post := &entity.Post{
		ID:      postID,
		Content: req.Content,
	}

	updatedPost, err := h.postUsecase.Update(c.Context(), userUUID, post)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to update post")
		return response.Error(c, http.StatusInternalServerError, err.Error(), 2004)
	}

	return response.Success(c, http.StatusOK, fiber.Map{
		"id":         updatedPost.ID,
		"thread_id":  updatedPost.ThreadID,
		"user_id":    updatedPost.UserID,
		"content":    updatedPost.Content,
		"created_at": updatedPost.CreatedAt,
		"updated_at": updatedPost.UpdatedAt,
	})
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete an existing post (requires authentication and ownership)
// @Tags posts
// @Security BearerAuth
// @Param id path string true "Post ID (UUID)"
// @Success 204 "Post deleted successfully"
// @Failure 400 {object} response.JSendResponse "Invalid post ID"
// @Failure 401 {object} response.JSendResponse "Unauthorized"
// @Failure 500 {object} response.JSendResponse "Internal server error"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse post ID from params")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid post ID"})
	}

	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		logger.GetLogger().Error().Msg("user_id not found or invalid in context")
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", 2001)
	}

	if err := h.postUsecase.Delete(c.Context(), userUUID, postID); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to delete post")
		return response.Error(c, http.StatusInternalServerError, err.Error(), 2005)
	}

	return c.Status(http.StatusNoContent).SendString("")
}

func RegisterPostRoutes(router fiber.Router, uc usecase.PostUsecase, jwtSecret string) {
	handler := NewPostHandler(uc)

	// Public routes
	router.Get("/posts/:id", handler.GetPost)
	router.Get("/users/:user_id/posts", handler.ListPostsByUser)

	// Protected routes
	protected := router.Group("/posts", middleware.AuthRequired(jwtSecret))
	protected.Post("/", handler.CreatePost)
	protected.Put("/:id", handler.UpdatePost)
	protected.Delete("/:id", handler.DeletePost)
}

type createPostRequest struct {
	ThreadID string `json:"thread_id" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

type updatePostRequest struct {
	Content string `json:"content" validate:"required"`
}
