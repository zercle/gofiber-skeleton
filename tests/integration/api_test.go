//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/delivery"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/repository"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/usecase"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/delivery"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/pkg/database"
	"github.com/zercle/gofiber-skeleton/tests/fixtures"
)

type APIIntegrationTestSuite struct {
	suite.Suite
	app             *fiber.App
	db              *database.Database
	config          *config.Config
	userRepo        repository.UserRepository
	postRepo        repository.PostRepository
	userUsecase     usecase.UserUsecase
	postUsecase     usecase.PostUsecase
	fixtures        *fixtures.UserFixtures
	postFixtures    *fixtures.PostFixtures
	requestFixtures *fixtures.RequestFixtures
	authToken       string
	testUser        *entity.User
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
	// Load test configuration
	cfg, err := config.LoadConfig("../testdata/test.env")
	suite.Require().NoError(err)

	// Use test database configuration
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}
	if cfg.Database.Port == "" {
		cfg.Database.Port = "5432"
	}
	if cfg.Database.Name == "" {
		cfg.Database.Name = "gofiber_skeleton_test"
	}
	if cfg.Database.User == "" {
		cfg.Database.User = "postgres"
	}
	if cfg.Database.Password == "" {
		cfg.Database.Password = "postgres"
	}

	// Configure JWT for testing
	cfg.JWT.Secret = "test-secret-key-for-integration-testing"
	cfg.JWT.Expiry = "24h"

	suite.config = cfg

	// Connect to test database
	db, err := database.NewDatabase(cfg.Database)
	suite.Require().NoError(err)

	suite.db = db
	suite.userRepo = repository.NewSQLCUserRepository(db)
	suite.postRepo = repository.NewSQLCPostRepository(db)

	// Initialize usecases
	suite.userUsecase = usecase.NewUserUsecase(suite.userRepo, cfg)
	suite.postUsecase = usecase.NewPostUsecase(suite.postRepo)

	// Setup Fiber app
	suite.app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	suite.setupRoutes()

	// Setup fixtures
	suite.fixtures = fixtures.NewUserFixtures()
	suite.postFixtures = fixtures.NewPostFixtures()
	suite.requestFixtures = fixtures.NewRequestFixtures()

	// Run migrations
	suite.runMigrations()

	// Create test user and get auth token
	suite.setupTestUser()
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.cleanupTestData()
		suite.db.Close()
	}
}

func (suite *APIIntegrationTestSuite) SetupTest() {
	// Clean up data before each test (except test user)
	suite.cleanupTestDataKeepTestUser()
}

func (suite *APIIntegrationTestSuite) TearDownTest() {
	// Clean up data after each test (except test user)
	suite.cleanupTestDataKeepTestUser()
}

func (suite *APIIntegrationTestSuite) runMigrations() {
	ctx := suite.app.Context()

	// Create users table
	_, err := suite.db.DB.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(255) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	suite.Require().NoError(err)

	// Create posts table
	_, err = suite.db.DB.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			title VARCHAR(255) NOT NULL,
			CONTENT TEXT NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	suite.Require().NoError(err)

	// Create indexes
	_, err = suite.db.DB.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	suite.Require().NoError(err)

	_, err = suite.db.DB.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id)`)
	suite.Require().NoError(err)
}

func (suite *APIIntegrationTestSuite) setupRoutes() {
	// User routes
	userHandler := delivery.NewUserHandler(suite.userUsecase)
	userGroup := suite.app.Group("/api/v1/users")
	userGroup.Post("/register", userHandler.Register)
	userGroup.Post("/login", userHandler.Login)
	userGroup.Get("/profile", userHandler.GetProfile)
	userGroup.Put("/profile", userHandler.UpdateProfile)
	userGroup.Get("/", userHandler.ListUsers)
	userGroup.Delete("/:id", userHandler.DeactivateUser)

	// Post routes
	postHandler := delivery.NewPostHandler(suite.postUsecase)
	postGroup := suite.app.Group("/api/v1/posts")
	postGroup.Get("/", postHandler.GetAllPosts)
	postGroup.Get("/published", postHandler.GetPublishedPosts)
	postGroup.Get("/search", postHandler.SearchPosts)
	postGroup.Post("/", postHandler.CreatePost)
	postGroup.Get("/:id", postHandler.GetPost)
	postGroup.Put("/:id", postHandler.UpdatePost)
	postGroup.Delete("/:id", postHandler.DeletePost)
	postGroup.Patch("/:id/publish", postHandler.PublishPost)
	postGroup.Patch("/:id/archive", postHandler.ArchivePost)
	postGroup.Patch("/:id/unpublish", postHandler.UnpublishPost)
	postGroup.Get("/user/:userId", postHandler.GetUserPosts)
	postGroup.Get("/user/:userId/stats", postHandler.GetUserPostStats)
}

func (suite *APIIntegrationTestSuite) setupTestUser() {
	// Create a test user for authentication
	user := entity.NewUser("testuser@example.com", "hashed_password", "Test User")
	err := suite.userRepo.Create(suite.app.Context(), user)
	suite.Require().NoError(err)

	suite.testUser = user

	// Login to get auth token
	loginReq := suite.requestFixtures.GetLoginRequest("valid")
	loginReq.Email = "testuser@example.com"

	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var loginResp struct {
		Token string       `json:"token"`
		User  *entity.User `json:"user"`
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	suite.Require().NoError(err)

	suite.authToken = loginResp.Token
}

func (suite *APIIntegrationTestSuite) cleanupTestData() {
	ctx := suite.app.Context()

	// Delete all posts first (due to foreign key constraint)
	_, err := suite.db.DB.Exec(ctx, "DELETE FROM posts")
	suite.Require().NoError(err)

	// Delete all users
	_, err = suite.db.DB.Exec(ctx, "DELETE FROM users")
	suite.Require().NoError(err)
}

func (suite *APIIntegrationTestSuite) cleanupTestDataKeepTestUser() {
	ctx := suite.app.Context()

	// Delete all posts
	_, err := suite.db.DB.Exec(ctx, "DELETE FROM posts")
	suite.Require().NoError(err)

	// Delete all users except test user
	_, err = suite.db.DB.Exec(ctx, "DELETE FROM users WHERE email != 'testuser@example.com'")
	suite.Require().NoError(err)
}

// Helper method to make authenticated requests
func (suite *APIIntegrationTestSuite) makeAuthenticatedRequest(method, path string, body []byte) (*http.Response, error) {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+suite.authToken)
	return suite.app.Test(req)
}

// Test User API Endpoints
func (suite *APIIntegrationTestSuite) TestUserAPI_Register() {
	tests := []struct {
		name           string
		requestFixture string
		expectedStatus int
		expectError    string
	}{
		{
			name:           "Valid registration",
			requestFixture: "valid",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Missing email",
			requestFixture: "missing_email",
			expectedStatus: http.StatusBadRequest,
			expectError:    "email is required",
		},
		{
			name:           "Missing password",
			requestFixture: "missing_password",
			expectedStatus: http.StatusBadRequest,
			expectError:    "password is required",
		},
		{
			name:           "Missing name",
			requestFixture: "missing_name",
			expectedStatus: http.StatusBadRequest,
			expectError:    "full_name is required",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req := suite.requestFixtures.GetCreateUserRequest(tt.requestFixture)
			reqBody, _ := json.Marshal(req)

			resp, err := suite.app.Test(httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewReader(reqBody)))
			suite.Require().NoError(err)
			suite.Equal(tt.expectedStatus, resp.StatusCode)

			if tt.expectError != "" {
				var errorResp struct {
					Error string `json:"error"`
				}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				suite.Require().NoError(err)
				suite.Contains(errorResp.Error, tt.expectError)
			}
		})
	}
}

func (suite *APIIntegrationTestSuite) TestUserAPI_Login() {
	// Create a user for login testing
	user := entity.NewUser("login@example.com", "hashed_password", "Login User")
	err := suite.userRepo.Create(suite.app.Context(), user)
	suite.Require().NoError(err)

	req := entity.LoginRequest{
		Email:    "login@example.com",
		Password: "password123", // This would need to be properly hashed in real implementation
	}

	reqBody, _ := json.Marshal(req)
	resp, err := suite.app.Test(httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewReader(reqBody)))
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var loginResp struct {
		Token string       `json:"token"`
		User  *entity.User `json:"user"`
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	suite.Require().NoError(err)
	suite.NotEmpty(loginResp.Token)
	suite.NotNil(loginResp.User)
	suite.Equal(user.Email, loginResp.User.Email)
}

func (suite *APIIntegrationTestSuite) TestUserAPI_GetProfile() {
	req := httptest.NewRequest("GET", "/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var profileResp *entity.User
	err = json.NewDecoder(resp.Body).Decode(&profileResp)
	suite.Require().NoError(err)
	suite.Equal(suite.testUser.Email, profileResp.Email)
	suite.Equal(suite.testUser.FullName, profileResp.FullName)
}

func (suite *APIIntegrationTestSuite) TestUserAPI_UpdateProfile() {
	updateReq := struct {
		Email    *string `json:"email,omitempty"`
		FullName *string `json:"full_name,omitempty"`
	}{
		Email:    stringPtr("updated@example.com"),
		FullName: stringPtr("Updated Name"),
	}

	reqBody, _ := json.Marshal(updateReq)
	resp, err := suite.makeAuthenticatedRequest("PUT", "/api/v1/users/profile", reqBody)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var profileResp *entity.User
	err = json.NewDecoder(resp.Body).Decode(&profileResp)
	suite.Require().NoError(err)
	suite.Equal("updated@example.com", profileResp.Email)
	suite.Equal("Updated Name", profileResp.FullName)
}

// Test Post API Endpoints
func (suite *APIIntegrationTestSuite) TestPostAPI_CreatePost() {
	tests := []struct {
		name           string
		requestFixture string
		expectedStatus int
		expectError    string
	}{
		{
			name:           "Valid draft post",
			requestFixture: "valid_draft",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Valid published post",
			requestFixture: "valid_published",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Missing title",
			requestFixture: "missing_title",
			expectedStatus: http.StatusBadRequest,
			expectError:    "title is required",
		},
		{
			name:           "Missing content",
			requestFixture: "missing_content",
			expectedStatus: http.StatusBadRequest,
			expectError:    "content is required",
		},
		{
			name:           "Invalid status",
			requestFixture: "invalid_status",
			expectedStatus: http.StatusBadRequest,
			expectError:    "invalid status",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req := suite.requestFixtures.GetCreatePostRequest(tt.requestFixture)
			reqBody, _ := json.Marshal(req)

			resp, err := suite.makeAuthenticatedRequest("POST", "/api/v1/posts", reqBody)
			suite.Require().NoError(err)
			suite.Equal(tt.expectedStatus, resp.StatusCode)

			if tt.expectError != "" {
				var errorResp struct {
					Error string `json:"error"`
				}
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				suite.Require().NoError(err)
				suite.Contains(errorResp.Error, tt.expectError)
			}
		})
	}
}

func (suite *APIIntegrationTestSuite) TestPostAPI_GetAndUpdatePost() {
	// Create a post first
	createReq := suite.requestFixtures.GetCreatePostRequest("valid_draft")
	reqBody, _ := json.Marshal(createReq)

	resp, err := suite.makeAuthenticatedRequest("POST", "/api/v1/posts", reqBody)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var createdPost *entity.Post
	err = json.NewDecoder(resp.Body).Decode(&createdPost)
	suite.Require().NoError(err)
	suite.NotEmpty(createdPost.ID)

	// Get the post
	resp, err = suite.app.Test(httptest.NewRequest("GET", "/api/v1/posts/"+createdPost.ID.String(), nil))
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var retrievedPost *entity.Post
	err = json.NewDecoder(resp.Body).Decode(&retrievedPost)
	suite.Require().NoError(err)
	suite.Equal(createdPost.ID, retrievedPost.ID)

	// Update the post
	updateReq := suite.requestFixtures.GetUpdatePostRequest("update_all")
	reqBody, _ = json.Marshal(updateReq)

	resp, err = suite.makeAuthenticatedRequest("PUT", "/api/v1/posts/"+createdPost.ID.String(), reqBody)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var updatedPost *entity.Post
	err = json.NewDecoder(resp.Body).Decode(&updatedPost)
	suite.Require().NoError(err)
	suite.Equal(*updateReq.Title, updatedPost.Title)
	suite.Equal(*updateReq.Content, updatedPost.Content)
	suite.Equal(*updateReq.Status, updatedPost.Status)
}

func (suite *APIIntegrationTestSuite) TestPostAPI_DeletePost() {
	// Create a post first
	createReq := suite.requestFixtures.GetCreatePostRequest("valid_draft")
	reqBody, _ := json.Marshal(createReq)

	resp, err := suite.makeAuthenticatedRequest("POST", "/api/v1/posts", reqBody)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var createdPost *entity.Post
	err = json.NewDecoder(resp.Body).Decode(&createdPost)
	suite.Require().NoError(err)

	// Delete the post
	resp, err = suite.makeAuthenticatedRequest("DELETE", "/api/v1/posts/"+createdPost.ID.String(), nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// Verify deletion
	resp, err = suite.app.Test(httptest.NewRequest("GET", "/api/v1/posts/"+createdPost.ID.String(), nil))
	suite.Require().NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *APIIntegrationTestSuite) TestPostAPI_StatusOperations() {
	// Create a post first
	createReq := suite.requestFixtures.GetCreatePostRequest("valid_draft")
	reqBody, _ := json.Marshal(createReq)

	resp, err := suite.makeAuthenticatedRequest("POST", "/api/v1/posts", reqBody)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var createdPost *entity.Post
	err = json.NewDecoder(resp.Body).Decode(&createdPost)
	suite.Require().NoError(err)
	suite.Equal(entity.PostStatusDraft, createdPost.Status)

	// Publish the post
	resp, err = suite.makeAuthenticatedRequest("PATCH", "/api/v1/posts/"+createdPost.ID.String()+"/publish", nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// Archive the post
	resp, err = suite.makeAuthenticatedRequest("PATCH", "/api/v1/posts/"+createdPost.ID.String()+"/archive", nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// Unpublish the post
	resp, err = suite.makeAuthenticatedRequest("PATCH", "/api/v1/posts/"+createdPost.ID.String()+"/unpublish", nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *APIIntegrationTestSuite) TestPostAPI_PublicEndpoints() {
	// Create some posts for testing
	user := entity.NewUser("postauthor@example.com", "hash", "Post Author")
	err := suite.userRepo.Create(suite.app.Context(), user)
	suite.Require().NoError(err)

	posts := []*entity.Post{
		entity.NewPost("Published Post 1", "Content 1", user.ID),
		entity.NewPost("Published Post 2", "Content 2", user.ID),
		entity.NewPost("Draft Post", "Draft content", user.ID),
	}

	posts[0].Status = entity.PostStatusPublished
	posts[1].Status = entity.PostStatusPublished
	posts[2].Status = entity.PostStatusDraft

	for _, post := range posts {
		err = suite.postRepo.Create(suite.app.Context(), post)
		suite.Require().NoError(err)
	}

	// Test get all posts
	resp, err := suite.app.Test(httptest.NewRequest("GET", "/api/v1/posts", nil))
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var allPosts []*entity.PostWithAuthor
	err = json.NewDecoder(resp.Body).Decode(&allPosts)
	suite.Require().NoError(err)
	suite.Len(allPosts, 3) // All posts returned for this endpoint

	// Test get published posts
	resp, err = suite.app.Test(httptest.NewRequest("GET", "/api/v1/posts/published", nil))
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var publishedPosts []*entity.PostWithAuthor
	err = json.NewDecoder(resp.Body).Decode(&publishedPosts)
	suite.Require().NoError(err)
	suite.Len(publishedPosts, 2) // Only published posts

	// Test search posts
	resp, err = suite.app.Test(httptest.NewRequest("GET", "/api/v1/posts/search?search=Published", nil))
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var searchResults []*entity.PostWithAuthor
	err = json.NewDecoder(resp.Body).Decode(&searchResults)
	suite.Require().NoError(err)
	suite.Len(searchResults, 2) // Posts with "Published" in title
}

func (suite *APIIntegrationTestSuite) TestPostAPI_UserSpecificEndpoints() {
	// Get user posts
	resp, err := suite.makeAuthenticatedRequest("GET", "/api/v1/posts/user/"+suite.testUser.ID.String(), nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var userPosts []*entity.Post
	err = json.NewDecoder(resp.Body).Decode(&userPosts)
	suite.Require().NoError(err)
	suite.Len(userPosts, 0) // No posts yet

	// Get user post stats
	resp, err = suite.makeAuthenticatedRequest("GET", "/api/v1/posts/user/"+suite.testUser.ID.String()+"/stats", nil)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var stats *entity.PostStats
	err = json.NewDecoder(resp.Body).Decode(&stats)
	suite.Require().NoError(err)
	suite.Equal(0, stats.TotalPosts)
}

func TestAPIIntegrationSuite(t *testing.T) {
	// Skip if integration tests are disabled
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("Skipping integration tests. Set INTEGRATION_TESTS=1 to run.")
	}

	suite.Run(t, new(APIIntegrationTestSuite))
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
