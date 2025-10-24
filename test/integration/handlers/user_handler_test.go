package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/template-go-fiber/internal/domains"
	"github.com/zercle/template-go-fiber/internal/handlers"
	"github.com/zercle/template-go-fiber/internal/usecases"
)

// MockUserRepository is a simple mock for testing
type MockUserRepository struct {
	users map[string]*domains.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domains.User),
	}
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domains.User, error) {
	return m.users[id], nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domains.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *domains.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domains.User) error {
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if user, ok := m.users[id]; ok {
		now := time.Now()
		user.DeletedAt = &now
	}
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int32) ([]*domains.User, error) {
	var users []*domains.User
	for _, user := range m.users {
		if user.DeletedAt == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	count := 0
	for _, user := range m.users {
		if user.DeletedAt == nil {
			count++
		}
	}
	return int64(count), nil
}

func setupTestApp(repo *MockUserRepository) *fiber.App {
	app := fiber.New()
	userUsecase := usecases.NewUserUsecase(repo)
	userHandler := handlers.NewUserHandler(userUsecase)

	api := app.Group("/api")
	api.Post("/users/register", userHandler.RegisterUser)
	api.Get("/users", userHandler.ListUsers)
	api.Get("/users/email", userHandler.GetUserByEmail)
	api.Get("/users/:id", userHandler.GetUser)
	api.Put("/users/:id", userHandler.UpdateUser)
	api.Delete("/users/:id", userHandler.DeleteUser)

	return app
}

func TestRegisterUserEndpoint(t *testing.T) {
	mockRepo := NewMockUserRepository()
	app := setupTestApp(mockRepo)
	defer func() {
		_ = app.Shutdown()
	}()

	body := `{"email":"test@example.com","password":"password123","first_name":"John","last_name":"Doe"}`
	req := httptest.NewRequest("POST", "/api/users/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &result)

	if success, ok := result["success"].(bool); !ok || !success {
		t.Error("expected success to be true")
	}
}

func TestGetUserEndpoint(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.users["123"] = &domains.User{ID: "123", Email: "user@example.com", IsActive: true}

	app := setupTestApp(mockRepo)
	defer func() {
		_ = app.Shutdown()
	}()

	req := httptest.NewRequest("GET", "/api/users/123", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetUserEndpoint_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	app := setupTestApp(mockRepo)
	defer func() {
		_ = app.Shutdown()
	}()

	req := httptest.NewRequest("GET", "/api/users/nonexistent", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestListUsersEndpoint(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.users["1"] = &domains.User{ID: "1", Email: "user1@example.com", IsActive: true}
	mockRepo.users["2"] = &domains.User{ID: "2", Email: "user2@example.com", IsActive: true}

	app := setupTestApp(mockRepo)
	defer func() {
		_ = app.Shutdown()
	}()

	req := httptest.NewRequest("GET", "/api/users?limit=10&offset=0", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDeleteUserEndpoint(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.users["123"] = &domains.User{ID: "123", Email: "user@example.com"}

	app := setupTestApp(mockRepo)
	defer func() {
		_ = app.Shutdown()
	}()

	req := httptest.NewRequest("DELETE", "/api/users/123", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestUpdateUserEndpoint(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.users["123"] = &domains.User{ID: "123", Email: "old@example.com"}

	app := setupTestApp(mockRepo)
	defer func() {
		_ = app.Shutdown()
	}()

	body := `{"email":"new@example.com","first_name":"Updated"}`
	req := httptest.NewRequest("PUT", "/api/users/123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
