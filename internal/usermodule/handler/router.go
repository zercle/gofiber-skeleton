package userhandler

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/zercle/gofiber-skeleton/docs" // Import for swagger docs
)

func SetupRoutes(router fiber.Router, jwtHandler fiber.Handler, handler *UserHandler) {
	setupPublicRoutes(router, handler)
	setupProtectedRoutes(router, jwtHandler, handler)
}

// SetupPublicRoutes initializes public user routes
func setupPublicRoutes(router fiber.Router, handler *UserHandler) {
	// User authentication & authorization
	userAPI := router.Group("/users")

	// Register godoc
	// @Summary Register a new user
	// @Description Register a new user with username and password
	// @Tags User
	// @Accept json
	// @Produce json
	// @Param request body RegisterRequest true "Register Request"
	// @Success 200 {object} usermodule.UserLoginRes "Success"
	// @Failure 400 {object} jsend.ErrorResponse "Bad Request"
	// @Failure 500 {object} jsend.ErrorResponse "Internal Server Error"
	// @Router /register [post]
	userAPI.Post("/register", handler.Register)

	// Login godoc
	// @Summary Login a user
	// @Description Login a user with username and password
	// @Tags User
	// @Accept json
	// @Produce json
	// @Param request body LoginRequest true "Login Request"
	// @Success 200 {object} usermodule.UserLoginRes "Success"
	// @Failure 400 {object} jsend.ErrorResponse "Bad Request"
	// @Failure 500 {object} jsend.ErrorResponse "Internal Server Error"
	// @Router /login [post]
	userAPI.Post("/login", handler.Login)
}

// SetupProtectedRoutes initializes protected user routes
func setupProtectedRoutes(router fiber.Router, jwtHandler fiber.Handler, handler *UserHandler) {

	userAPI := router.Group("/users")
	userAPI.Use(jwtHandler)

	// Protected user routes
	// @Summary Get user by ID
	// @Description Get a user by their ID
	// @Tags User
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param id path string true "User ID"
	// @Success 200 {object} usermodule.UserLoginRes "Success"
	// @Failure 400 {object} jsend.ErrorResponse "Bad Request"
	// @Failure 401 {object} jsend.ErrorResponse "Unauthorized"
	// @Failure 404 {object} jsend.ErrorResponse "Not Found"
	// @Failure 500 {object} jsend.ErrorResponse "Internal Server Error"
	// @Router /users/{id} [get]
	userAPI.Get("/users/:id", handler.GetByID)

	// @Summary Update user role
	// @Description Update a user's role by ID (Admin only)
	// @Tags User
	// @Accept json
	// @Produce json
	// @Security ApiKeyAuth
	// @Param id path string true "User ID"
	// @Param request body userhandler.UpdateRoleRequest true "Update Role Request"
	// @Success 200 {object} jsend.JSendResponse "Success"
	// @Failure 400 {object} jsend.ErrorResponse "Bad Request"
	// @Failure 401 {object} jsend.ErrorResponse "Unauthorized"
	// @Failure 403 {object} jsend.ErrorResponse "Forbidden"
	// @Failure 404 {object} jsend.ErrorResponse "Not Found"
	// @Failure 500 {object} jsend.ErrorResponse "Internal Server Error"
	// @Router /users/{id}/role [put]
	userAPI.Put("/users/:id/role", handler.UpdateRole)
}
