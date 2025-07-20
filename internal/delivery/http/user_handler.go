package http

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-skeleton/internal/auth"
	"gofiber-skeleton/internal/configs"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Registers a new user with a username and password.
// @Tags users
// @Accept json
// @Produce json
// @Param user body registerRequest true "User registration details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /register [post]
func (h *Handler) RegisterUser(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "data": err.Error()})
	}

	user, err := h.userUsecase.Register(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	cfg, err := configs.LoadConfig()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to load config"})
	}

	token, err := auth.GenerateToken(user.ID, cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user, "token": token}})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token.
// @Tags users
// @Accept json
// @Produce json
// @Param user body registerRequest true "User login details"
// @Success 200 {object} map[string]interface{} "User logged in successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "data": err.Error()})
	}

	user, err := h.userUsecase.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "data": "Invalid credentials"})
	}

	cfg, err := configs.LoadConfig()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to load config"})
	}

	token, err := auth.GenerateToken(user.ID, cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"status": "success", "data": fiber.Map{"token": token}})
}
