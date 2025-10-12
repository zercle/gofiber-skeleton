package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/pkg/logger"
)

// ApiResponse represents a standard API response
type ApiResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool        `json:"success"`
	Error     string      `json:"error"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// TodoItem represents a todo item for the sample API
type TodoItem struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIHandler handles sample API routes
type APIHandler struct {
	todos   []TodoItem
	nextID  int
}

// NewAPIHandler creates a new API handler instance
func NewAPIHandler() *APIHandler {
	todos := []TodoItem{
		{
			ID:        1,
			Title:     "Learn Go Fiber",
			Completed: false,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        2,
			Title:     "Build production-ready backend",
			Completed: true,
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-6 * time.Hour),
		},
	}

	return &APIHandler{
		todos:  todos,
		nextID: 3,
	}
}

// GetTodos returns all todo items
func (h *APIHandler) GetTodos(c *fiber.Ctx) error {
	logger.Info("Getting all todos", "count", len(h.todos))

	response := ApiResponse{
		Success:   true,
		Data:      h.todos,
		Message:   "Todos retrieved successfully",
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}

	return c.JSON(response)
}

// GetTodo returns a single todo item by ID
func (h *APIHandler) GetTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Error("Invalid todo ID", "id", c.Params("id"), "error", err)
		return h.sendError(c, fiber.StatusBadRequest, "Invalid todo ID", err)
	}

	for _, todo := range h.todos {
		if todo.ID == id {
			logger.Info("Todo retrieved successfully", "id", id, "title", todo.Title)
			response := ApiResponse{
				Success:   true,
				Data:      todo,
				Message:   fmt.Sprintf("Todo with ID %d retrieved successfully", id),
				Timestamp: time.Now(),
				RequestID: getRequestID(c),
			}
			return c.JSON(response)
		}
	}

	logger.Warn("Todo not found", "id", id)
	return h.sendError(c, fiber.StatusNotFound, fmt.Sprintf("Todo with ID %d not found", id), nil)
}

// CreateTodo creates a new todo item
func (h *APIHandler) CreateTodo(c *fiber.Ctx) error {
	var requestBody struct {
		Title string `json:"title"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Error("Failed to parse request body", "error", err)
		return h.sendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if requestBody.Title == "" {
		logger.Error("Empty todo title")
		return h.sendError(c, fiber.StatusBadRequest, "Title is required", nil)
	}

	now := time.Now()
	newTodo := TodoItem{
		ID:        h.nextID,
		Title:     requestBody.Title,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	h.todos = append(h.todos, newTodo)
	h.nextID++

	logger.Info("Todo created successfully", "id", newTodo.ID, "title", newTodo.Title)

	response := ApiResponse{
		Success:   true,
		Data:      newTodo,
		Message:   "Todo created successfully",
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateTodo updates an existing todo item
func (h *APIHandler) UpdateTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Error("Invalid todo ID", "id", c.Params("id"), "error", err)
		return h.sendError(c, fiber.StatusBadRequest, "Invalid todo ID", err)
	}

	var requestBody struct {
		Title     *string `json:"title,omitempty"`
		Completed *bool   `json:"completed,omitempty"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Error("Failed to parse request body", "error", err)
		return h.sendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	for i, todo := range h.todos {
		if todo.ID == id {
			if requestBody.Title != nil {
				h.todos[i].Title = *requestBody.Title
			}
			if requestBody.Completed != nil {
				h.todos[i].Completed = *requestBody.Completed
			}
			h.todos[i].UpdatedAt = time.Now()

			logger.Info("Todo updated successfully", "id", id, "title", h.todos[i].Title)

			response := ApiResponse{
				Success:   true,
				Data:      h.todos[i],
				Message:   "Todo updated successfully",
				Timestamp: time.Now(),
				RequestID: getRequestID(c),
			}
			return c.JSON(response)
		}
	}

	logger.Warn("Todo not found for update", "id", id)
	return h.sendError(c, fiber.StatusNotFound, fmt.Sprintf("Todo with ID %d not found", id), nil)
}

// DeleteTodo deletes a todo item
func (h *APIHandler) DeleteTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Error("Invalid todo ID", "id", c.Params("id"), "error", err)
		return h.sendError(c, fiber.StatusBadRequest, "Invalid todo ID", err)
	}

	for i, todo := range h.todos {
		if todo.ID == id {
			h.todos = append(h.todos[:i], h.todos[i+1:]...)

			logger.Info("Todo deleted successfully", "id", id, "title", todo.Title)

			response := ApiResponse{
				Success:   true,
				Message:   fmt.Sprintf("Todo with ID %d deleted successfully", id),
				Timestamp: time.Now(),
				RequestID: getRequestID(c),
			}
			return c.JSON(response)
		}
	}

	logger.Warn("Todo not found for deletion", "id", id)
	return h.sendError(c, fiber.StatusNotFound, fmt.Sprintf("Todo with ID %d not found", id), nil)
}

// GetStats returns API statistics
func (h *APIHandler) GetStats(c *fiber.Ctx) error {
	total := len(h.todos)
	completed := 0

	for _, todo := range h.todos {
		if todo.Completed {
			completed++
		}
	}

	stats := map[string]interface{}{
		"total_todos":     total,
		"completed_todos": completed,
		"pending_todos":   total - completed,
		"completion_rate": float64(completed) / float64(total) * 100,
	}

	logger.Info("API stats retrieved", "total_todos", total, "completed", completed)

	response := ApiResponse{
		Success:   true,
		Data:      stats,
		Message:   "API statistics retrieved successfully",
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}

	return c.JSON(response)
}

// getRequestID safely retrieves the request ID from context
func getRequestID(c *fiber.Ctx) string {
	if requestID := c.Locals("requestID"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// sendError sends a standardized error response
func (h *APIHandler) sendError(c *fiber.Ctx, statusCode int, message string, details interface{}) error {
	response := ErrorResponse{
		Success:   false,
		Error:     message,
		Details:   details,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}

	return c.Status(statusCode).JSON(response)
}