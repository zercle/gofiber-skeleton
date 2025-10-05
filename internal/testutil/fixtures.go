package testutil

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zercle/gofiber-skeleton/internal/post/entity"
	userEntity "github.com/zercle/gofiber-skeleton/internal/user/entity"
)

// UserFixture creates a test user with default values
func UserFixture(t *testing.T, overrides ...func(*userEntity.User)) *userEntity.User {
	t.Helper()

	id, err := uuid.NewV7()
	require.NoError(t, err, "failed to generate UUID")
	user := &userEntity.User{
		ID:        id,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: TimeNow(),
		UpdatedAt: TimeNow(),
	}

	for _, override := range overrides {
		override(user)
	}

	return user
}

// PostFixture creates a test post with default values
func PostFixture(t *testing.T, overrides ...func(*entity.Post)) *entity.Post {
	t.Helper()

	id, err := uuid.NewV7()
	require.NoError(t, err, "failed to generate UUID for post ID")
	userID, err := uuid.NewV7()
	require.NoError(t, err, "failed to generate UUID for user ID")
	threadID, err := uuid.NewV7()
	require.NoError(t, err, "failed to generate UUID for thread ID")

	post := &entity.Post{
		ID:        id,
		ThreadID:  threadID,
		UserID:    userID,
		Content:   "Test post content",
		CreatedAt: TimeNow(),
		UpdatedAt: TimeNow(),
	}

	for _, override := range overrides {
		override(post)
	}

	return post
}

// MultiplePostFixtures creates multiple test posts
func MultiplePostFixtures(t *testing.T, count int, userID uuid.UUID) []*entity.Post {
	t.Helper()

	posts := make([]*entity.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = PostFixture(t, func(p *entity.Post) {
			p.UserID = userID
		})
	}

	return posts
}

// PasswordHash returns a pre-computed bcrypt hash for "password123"
// This avoids slow bcrypt operations in tests
func PasswordHash() string {
	// bcrypt hash of "password123" with cost 10
	return "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
}

// ValidJWTSecret returns a test JWT secret
func ValidJWTSecret() string {
	return "test-jwt-secret-key-for-testing-only"
}

// TimePtr returns a pointer to a time value
func TimePtr(t time.Time) *time.Time {
	return &t
}

// StringPtr returns a pointer to a string value
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int value
func IntPtr(i int) *int {
	return &i
}
