//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/repository"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	"github.com/zercle/gofiber-skeleton/pkg/database"
)

type DatabaseIntegrationTestSuite struct {
	suite.Suite
	db       *database.Database
	userRepo repository.UserRepository
	postRepo repository.PostRepository
	config   *config.Config
}

func (suite *DatabaseIntegrationTestSuite) SetupSuite() {
	// Load test configuration
	cfg, err := config.LoadConfig("../testdata/test.env")
	suite.Require().NoError(err)

	// Use test database
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

	suite.config = cfg

	// Connect to test database
	db, err := database.NewDatabase(cfg.Database)
	suite.Require().NoError(err)

	suite.db = db
	suite.userRepo = repository.NewSQLCUserRepository(db)
	suite.postRepo = repository.NewSQLCPostRepository(db)

	// Run migrations
	suite.runMigrations()
}

func (suite *DatabaseIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		// Clean up test data
		suite.cleanupTestData()
		suite.db.Close()
	}
}

func (suite *DatabaseIntegrationTestSuite) SetupTest() {
	// Clean up data before each test
	suite.cleanupTestData()
}

func (suite *DatabaseIntegrationTestSuite) TearDownTest() {
	// Clean up data after each test
	suite.cleanupTestData()
}

func (suite *DatabaseIntegrationTestSuite) runMigrations() {
	// This would run database migrations
	// For now, we'll create the necessary tables manually
	suite.createTestTables()
}

func (suite *DatabaseIntegrationTestSuite) createTestTables() {
	ctx := context.Background()

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

	_, err = suite.db.DB.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status)`)
	suite.Require().NoError(err)

	_, err = suite.db.DB.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at)`)
	suite.Require().NoError(err)
}

func (suite *DatabaseIntegrationTestSuite) cleanupTestData() {
	ctx := context.Background()

	// Delete all posts first (due to foreign key constraint)
	_, err := suite.db.DB.Exec(ctx, "DELETE FROM posts")
	if err != nil {
		fmt.Printf("Error cleaning posts: %v\n", err)
	}

	// Delete all users
	_, err = suite.db.DB.Exec(ctx, "DELETE FROM users")
	if err != nil {
		fmt.Printf("Error cleaning users: %v\n", err)
	}
}

// Test User Repository Integration
func (suite *DatabaseIntegrationTestSuite) TestUserRepository_CreateAndGet() {
	ctx := context.Background()

	// Create a user
	user := entity.NewUser("test@example.com", "hashed_password", "Test User")

	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)
	suite.NotEmpty(user.ID)

	// Get user by ID
	retrievedUser, err := suite.userRepo.GetByID(ctx, user.ID)
	suite.Require().NoError(err)
	suite.Equal(user.ID, retrievedUser.ID)
	suite.Equal(user.Email, retrievedUser.Email)
	suite.Equal(user.FullName, retrievedUser.FullName)
	suite.True(retrievedUser.IsActive)

	// Get user by email
	retrievedUserByEmail, err := suite.userRepo.GetByEmail(ctx, user.Email)
	suite.Require().NoError(err)
	suite.Equal(user.ID, retrievedUserByEmail.ID)
	suite.Equal(user.Email, retrievedUserByEmail.Email)
}

func (suite *DatabaseIntegrationTestSuite) TestUserRepository_Update() {
	ctx := context.Background()

	// Create a user
	user := entity.NewUser("test@example.com", "hashed_password", "Test User")
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	// Update user
	user.FullName = "Updated User"
	user.UpdatedAt = time.Now()

	err = suite.userRepo.Update(ctx, user)
	suite.Require().NoError(err)

	// Verify update
	retrievedUser, err := suite.userRepo.GetByID(ctx, user.ID)
	suite.Require().NoError(err)
	suite.Equal("Updated User", retrievedUser.FullName)
}

func (suite *DatabaseIntegrationTestSuite) TestUserRepository_EmailExists() {
	ctx := context.Background()

	// Create a user
	user := entity.NewUser("test@example.com", "hashed_password", "Test User")
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	// Check if email exists
	exists, err := suite.userRepo.EmailExists(ctx, "test@example.com")
	suite.Require().NoError(err)
	suite.True(exists)

	// Check non-existent email
	exists, err = suite.userRepo.EmailExists(ctx, "nonexistent@example.com")
	suite.Require().NoError(err)
	suite.False(exists)
}

func (suite *DatabaseIntegrationTestSuite) TestUserRepository_ListAndCount() {
	ctx := context.Background()

	// Create multiple users
	users := []*entity.User{
		entity.NewUser("user1@example.com", "hash1", "User One"),
		entity.NewUser("user2@example.com", "hash2", "User Two"),
		entity.NewUser("user3@example.com", "hash3", "User Three"),
	}

	for _, user := range users {
		err := suite.userRepo.Create(ctx, user)
		suite.Require().NoError(err)
	}

	// Test count
	count, err := suite.userRepo.Count(ctx)
	suite.Require().NoError(err)
	suite.Equal(3, count)

	// Test list
	listedUsers, err := suite.userRepo.List(ctx, 10, 0)
	suite.Require().NoError(err)
	suite.Len(listedUsers, 3)

	// Test pagination
	paginatedUsers, err := suite.userRepo.List(ctx, 2, 0)
	suite.Require().NoError(err)
	suite.Len(paginatedUsers, 2)

	paginatedUsers, err = suite.userRepo.List(ctx, 2, 2)
	suite.Require().NoError(err)
	suite.Len(paginatedUsers, 1)
}

// Test Post Repository Integration
func (suite *DatabaseIntegrationTestSuite) TestPostRepository_CreateAndGet() {
	ctx := context.Background()

	// Create a user first
	user := entity.NewUser("test@example.com", "hashed_password", "Test User")
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	// Create a post
	post := entity.NewPost("Test Post", "This is test content", user.ID)

	err = suite.postRepo.Create(ctx, post)
	suite.Require().NoError(err)
	suite.NotEmpty(post.ID)

	// Get post by ID
	retrievedPost, err := suite.postRepo.GetByID(ctx, post.ID)
	suite.Require().NoError(err)
	suite.Equal(post.ID, retrievedPost.ID)
	suite.Equal(post.Title, retrievedPost.Title)
	suite.Equal(post.Content, retrievedPost.Content)
	suite.Equal(post.UserID, retrievedPost.UserID)
	suite.Equal(entity.PostStatusDraft, retrievedPost.Status)
}

func (suite *DatabaseIntegrationTestSuite) TestPostRepository_UserOperations() {
	ctx := context.Background()

	// Create a user
	user := entity.NewUser("test@example.com", "hashed_password", "Test User")
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	// Create multiple posts
	posts := []*entity.Post{
		entity.NewPost("Draft Post 1", "Content 1", user.ID),
		entity.NewPost("Draft Post 2", "Content 2", user.ID),
		entity.NewPost("Published Post 1", "Content 3", user.ID),
	}

	// Set one post as published
	posts[2].Status = entity.PostStatusPublished

	for _, post := range posts {
		err = suite.postRepo.Create(ctx, post)
		suite.Require().NoError(err)
	}

	// Test get all user posts
	userPosts, err := suite.postRepo.GetByUserID(ctx, user.ID)
	suite.Require().NoError(err)
	suite.Len(userPosts, 3)

	// Test get published posts
	publishedPosts, err := suite.postRepo.GetPublishedPostsByUserID(ctx, user.ID)
	suite.Require().NoError(err)
	suite.Len(publishedPosts, 1)
	suite.Equal(entity.PostStatusPublished, publishedPosts[0].Status)

	// Test get draft posts
	draftPosts, err := suite.postRepo.GetDraftPostsByUserID(ctx, user.ID)
	suite.Require().NoError(err)
	suite.Len(draftPosts, 2)

	for _, post := range draftPosts {
		suite.Equal(entity.PostStatusDraft, post.Status)
	}
}

func (suite *DatabaseIntegrationTestSuite) TestPostRepository_Authorization() {
	ctx := context.Background()

	// Create two users
	user1 := entity.NewUser("user1@example.com", "hash1", "User One")
	user2 := entity.NewUser("user2@example.com", "hash2", "User Two")

	err := suite.userRepo.Create(ctx, user1)
	suite.Require().NoError(err)
	err = suite.userRepo.Create(ctx, user2)
	suite.Require().NoError(err)

	// Create a post for user1
	post := entity.NewPost("User1 Post", "Content for user1", user1.ID)
	err = suite.postRepo.Create(ctx, post)
	suite.Require().NoError(err)

	// Test ownership check
	isOwner, err := suite.postRepo.IsOwner(ctx, post.ID, user1.ID)
	suite.Require().NoError(err)
	suite.True(isOwner)

	isOwner, err = suite.postRepo.IsOwner(ctx, post.ID, user2.ID)
	suite.Require().NoError(err)
	suite.False(isOwner)

	// Test existence check
	exists, err := suite.postRepo.Exists(ctx, post.ID)
	suite.Require().NoError(err)
	suite.True(exists)

	exists, err = suite.postRepo.Exists(ctx, uuid.New())
	suite.Require().NoError(err)
	suite.False(exists)
}

func (suite *DatabaseIntegrationTestSuite) TestPostRepository_UpdateDelete() {
	ctx := context.Background()

	// Create a user
	user := entity.NewUser("test@example.com", "hashed_password", "Test User")
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	// Create a post
	post := entity.NewPost("Original Title", "Original content", user.ID)
	err = suite.postRepo.Create(ctx, post)
	suite.Require().NoError(err)

	// Update post
	post.Title = "Updated Title"
	post.Content = "Updated content"
	post.Status = entity.PostStatusPublished
	post.UpdatedAt = time.Now()

	err = suite.postRepo.Update(ctx, post)
	suite.Require().NoError(err)

	// Verify update
	retrievedPost, err := suite.postRepo.GetByID(ctx, post.ID)
	suite.Require().NoError(err)
	suite.Equal("Updated Title", retrievedPost.Title)
	suite.Equal("Updated content", retrievedPost.Content)
	suite.Equal(entity.PostStatusPublished, retrievedPost.Status)

	// Delete post
	err = suite.postRepo.Delete(ctx, post.ID)
	suite.Require().NoError(err)

	// Verify deletion
	_, err = suite.postRepo.GetByID(ctx, post.ID)
	suite.Error(err)
	suite.Equal(entity.ErrPostNotFound, err)
}

func (suite *DatabaseIntegrationTestSuite) TestPostRepository_SearchAndStats() {
	ctx := context.Background()

	// Create users
	author1 := entity.NewUser("author1@example.com", "hash1", "Author One")
	author2 := entity.NewUser("author2@example.com", "hash2", "Author Two")

	err := suite.userRepo.Create(ctx, author1)
	suite.Require().NoError(err)
	err = suite.userRepo.Create(ctx, author2)
	suite.Require().NoError(err)

	// Create posts with different content
	posts := []*entity.Post{
		entity.NewPost("Test Post 1", "This is about testing", author1.ID),
		entity.NewPost("Tutorial Post", "This is a tutorial", author1.ID),
		entity.NewPost("Another Test", "More test content", author2.ID),
	}

	// Set statuses
	posts[0].Status = entity.PostStatusPublished
	posts[1].Status = entity.PostStatusPublished
	posts[2].Status = entity.PostStatusDraft

	for _, post := range posts {
		err = suite.postRepo.Create(ctx, post)
		suite.Require().NoError(err)
	}

	// Test search
	searchResults, err := suite.postRepo.SearchPosts(ctx, "test", entity.PostStatusPublished, 10, 0)
	suite.Require().NoError(err)
	suite.Len(searchResults, 1) // Only published posts with "test"
	suite.Equal("Test Post 1", searchResults[0].Title)

	// Test user stats
	stats, err := suite.postRepo.GetUserPostStats(ctx, author1.ID)
	suite.Require().NoError(err)
	suite.Equal(2, stats.TotalPosts)
	suite.Equal(2, stats.PublishedPosts)
	suite.Equal(0, stats.DraftPosts)
	suite.NotNil(stats.LastPostDate)

	// Test status counts
	statusCounts, err := suite.postRepo.CountPostsByStatus(ctx)
	suite.Require().NoError(err)
	suite.Len(statusCounts, 2) // published and draft

	var publishedCount, draftCount int
	for _, count := range statusCounts {
		if count.Status == "published" {
			publishedCount = count.Count
		} else if count.Status == "draft" {
			draftCount = count.Count
		}
	}
	suite.Equal(2, publishedCount)
	suite.Equal(1, draftCount)
}

// Test Runner
func TestDatabaseIntegrationSuite(t *testing.T) {
	// Skip if integration tests are disabled
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("Skipping integration tests. Set INTEGRATION_TESTS=1 to run.")
	}

	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
