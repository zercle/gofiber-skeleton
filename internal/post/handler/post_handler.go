package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/post/usecase"
	"github.com/zercle/gofiber-skeleton/internal/user/middleware"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
}

func NewPostHandler(postUsecase usecase.PostUsecase) *PostHandler {
	return &PostHandler{postUsecase: postUsecase}
}

func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	req := new(createPostRequest)
	if err := c.BodyParser(req); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse create post request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// TODO: Add validation using a library like go-playground/validator

	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse user ID from context")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}
	
	threadID, err := uuid.Parse(req.ThreadID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse thread ID from request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid thread ID"})
	}

	post, err := h.postUsecase.Create(c.Context(), userUUID, threadID, req.Content)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to create post")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create post"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":         post.ID,
		"thread_id":  post.ThreadID,
		"user_id":    post.UserID,
		"content":    post.Content,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	})
}

func (h *PostHandler) GetPost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse post ID from params")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid post ID"})
	}

	post, err := h.postUsecase.Get(c.Context(), postID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to get post")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Post not found"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id":         post.ID,
		"thread_id":  post.ThreadID,
		"user_id":    post.UserID,
		"content":    post.Content,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	})
}

func (h *PostHandler) ListPostsByUser(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse user ID from params")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	posts, err := h.postUsecase.ListByUser(c.Context(), userID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to list posts by user")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve posts"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"posts": posts})
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse post ID from params")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid post ID"})
	}

	req := new(updatePostRequest)
	if err := c.BodyParser(req); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse update post request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// TODO: Add validation using a library like go-playground/validator

	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse user ID from context")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	post := &entity.Post{
		ID:      postID,
		Content: req.Content,
	}

	updatedPost, err := h.postUsecase.Update(c.Context(), userUUID, post)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to update post")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id":         updatedPost.ID,
		"thread_id":  updatedPost.ThreadID,
		"user_id":    updatedPost.UserID,
		"content":    updatedPost.Content,
		"created_at": updatedPost.CreatedAt,
		"updated_at": updatedPost.UpdatedAt,
	})
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse post ID from params")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid post ID"})
	}

	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse user ID from context")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	if err := h.postUsecase.Delete(c.Context(), userUUID, postID); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to delete post")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
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