package fixtures

import (
	"time"

	postentity "github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	userentity "github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/pkg/uuid"
	uuidv7 "github.com/google/uuid"
)

// UserFixtures provides predefined user entities for testing
type UserFixtures struct {
	users map[string]*userentity.DomainUser
}

// NewUserFixtures creates a new set of user fixtures
func NewUserFixtures() *UserFixtures {
	users := make(map[string]*userentity.DomainUser)

	// Create test users using deterministic UUIDv7 values
	users["admin"] = &userentity.DomainUser{
		ID:           uuid.GenerateTestUUIDv7FromString("admin-user"),
		Email:        "admin@example.com",
		PasswordHash: "hashed_admin_password",
		FullName:     "Admin User",
		IsActive:     true,
		CreatedAt:    time.Now().Add(-24 * time.Hour),
		UpdatedAt:    time.Now().Add(-1 * time.Hour),
	}

	users["regular"] = &userentity.DomainUser{
		ID:           uuid.GenerateTestUUIDv7FromString("regular-user"),
		Email:        "user@example.com",
		PasswordHash: "hashed_user_password",
		FullName:     "Regular User",
		IsActive:     true,
		CreatedAt:    time.Now().Add(-12 * time.Hour),
		UpdatedAt:    time.Now().Add(-30 * time.Minute),
	}

	users["inactive"] = &userentity.DomainUser{
		ID:           uuid.GenerateTestUUIDv7FromString("inactive-user"),
		Email:        "inactive@example.com",
		PasswordHash: "hashed_inactive_password",
		FullName:     "Inactive User",
		IsActive:     false,
		CreatedAt:    time.Now().Add(-48 * time.Hour),
		UpdatedAt:    time.Now().Add(-2 * time.Hour),
	}

	return &UserFixtures{users: users}
}

// Get returns a user fixture by name
func (f *UserFixtures) Get(name string) *userentity.DomainUser {
	if user, exists := f.users[name]; exists {
		// Return a copy to avoid modifying the fixture
		userCopy := *user
		return &userCopy
	}
	return nil
}

// GetAll returns all user fixtures
func (f *UserFixtures) GetAll() map[string]*userentity.DomainUser {
	users := make(map[string]*userentity.DomainUser)
	for name, user := range f.users {
		userCopy := *user
		users[name] = &userCopy
	}
	return users
}

// PostFixtures provides predefined post entities for testing
type PostFixtures struct {
	posts map[string]*postentity.DomainPost
}

// NewPostFixtures creates a new set of post fixtures
func NewPostFixtures() *PostFixtures {
	posts := make(map[string]*postentity.DomainPost)

	adminID := uuid.GenerateTestUUIDv7FromString("admin-user")
	userID := uuid.GenerateTestUUIDv7FromString("regular-user")

	// Create test posts using deterministic UUIDv7 values
	posts["admin_published"] = &postentity.DomainPost{
		ID:        uuid.GenerateTestUUIDv7FromString("admin-published-post"),
		Title:     "Admin Published Post",
		Content:   "This is a published post by admin",
		Status:    postentity.PostStatusPublished,
		UserID:    adminID,
		CreatedAt: time.Now().Add(-6 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	posts["admin_draft"] = &postentity.DomainPost{
		ID:        uuid.GenerateTestUUIDv7FromString("admin-draft-post"),
		Title:     "Admin Draft Post",
		Content:   "This is a draft post by admin",
		Status:    postentity.PostStatusDraft,
		UserID:    adminID,
		CreatedAt: time.Now().Add(-3 * time.Hour),
		UpdatedAt: time.Now().Add(-30 * time.Minute),
	}

	posts["user_published"] = &postentity.DomainPost{
		ID:        uuid.GenerateTestUUIDv7FromString("user-published-post"),
		Title:     "User Published Post",
		Content:   "This is a published post by regular user",
		Status:    postentity.PostStatusPublished,
		UserID:    userID,
		CreatedAt: time.Now().Add(-8 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}

	posts["user_draft"] = &postentity.DomainPost{
		ID:        uuid.GenerateTestUUIDv7FromString("user-draft-post"),
		Title:     "User Draft Post",
		Content:   "This is a draft post by regular user",
		Status:    postentity.PostStatusDraft,
		UserID:    userID,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-15 * time.Minute),
	}

	posts["archived"] = &postentity.DomainPost{
		ID:        uuid.GenerateTestUUIDv7FromString("archived-post"),
		Title:     "Archived Post",
		Content:   "This is an archived post",
		Status:    postentity.PostStatusArchived,
		UserID:    adminID,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-12 * time.Hour),
	}

	return &PostFixtures{posts: posts}
}

// Get returns a post fixture by name
func (f *PostFixtures) Get(name string) *postentity.DomainPost {
	if post, exists := f.posts[name]; exists {
		// Return a copy to avoid modifying the fixture
		postCopy := *post
		return &postCopy
	}
	return nil
}

// GetAll returns all post fixtures
func (f *PostFixtures) GetAll() map[string]*postentity.DomainPost {
	posts := make(map[string]*postentity.DomainPost)
	for name, post := range f.posts {
		postCopy := *post
		posts[name] = &postCopy
	}
	return posts
}

// GetByStatus returns posts filtered by status
func (f *PostFixtures) GetByStatus(status postentity.PostStatus) []*postentity.DomainPost {
	var result []*postentity.DomainPost
	for _, post := range f.posts {
		if post.Status == status {
			postCopy := *post
			result = append(result, &postCopy)
		}
	}
	return result
}

// GetByUserID returns posts filtered by user ID
func (f *PostFixtures) GetByUserID(userID uuidv7.UUID) []*postentity.DomainPost {
	var result []*postentity.DomainPost
	for _, post := range f.posts {
		if post.UserID == userID {
			postCopy := *post
			result = append(result, &postCopy)
		}
	}
	return result
}

// PostWithAuthorFixtures provides predefined post with author entities for testing
type PostWithAuthorFixtures struct {
	posts map[string]*postentity.PostWithAuthor
}

// NewPostWithAuthorFixtures creates a new set of post with author fixtures
func NewPostWithAuthorFixtures() *PostWithAuthorFixtures {
	posts := make(map[string]*postentity.PostWithAuthor)

	adminID := uuid.GenerateTestUUIDv7FromString("admin-user")
	userID := uuid.GenerateTestUUIDv7FromString("regular-user")

	// Create test posts with author info using deterministic UUIDv7 values
	posts["admin_post"] = &postentity.PostWithAuthor{
		ID:        uuid.GenerateTestUUIDv7FromString("admin-post-with-author"),
		Title:     "Admin Post with Author",
		Content:   "This post includes author information",
		Status:    postentity.PostStatusPublished,
		UserID:    adminID,
		FullName:  "Admin User",
		Email:     "admin@example.com",
		CreatedAt: time.Now().Add(-4 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	posts["user_post"] = &postentity.PostWithAuthor{
		ID:        uuid.GenerateTestUUIDv7FromString("user-post-with-author"),
		Title:     "User Post with Author",
		Content:   "This post includes author information",
		Status:    postentity.PostStatusPublished,
		UserID:    userID,
		FullName:  "Regular User",
		Email:     "user@example.com",
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-30 * time.Minute),
	}

	return &PostWithAuthorFixtures{posts: posts}
}

// Get returns a post with author fixture by name
func (f *PostWithAuthorFixtures) Get(name string) *postentity.PostWithAuthor {
	if post, exists := f.posts[name]; exists {
		postCopy := *post
		return &postCopy
	}
	return nil
}

// GetAll returns all post with author fixtures
func (f *PostWithAuthorFixtures) GetAll() map[string]*postentity.PostWithAuthor {
	posts := make(map[string]*postentity.PostWithAuthor)
	for name, post := range f.posts {
		postCopy := *post
		posts[name] = &postCopy
	}
	return posts
}

// RequestFixtures provides predefined request DTOs for testing
type RequestFixtures struct {
	createUserRequests map[string]*userentity.CreateUserRequest
	loginRequests      map[string]*userentity.LoginRequest
	createPostRequests map[string]*postentity.CreatePostRequest
	updatePostRequests map[string]*postentity.UpdatePostRequest
}

// NewRequestFixtures creates a new set of request fixtures
func NewRequestFixtures() *RequestFixtures {
	createUserRequests := map[string]*userentity.CreateUserRequest{
		"valid": {
			Email:    "newuser@example.com",
			Password: "password123",
			FullName: "New User",
		},
		"missing_email": {
			Email:    "",
			Password: "password123",
			FullName: "New User",
		},
		"missing_password": {
			Email:    "newuser@example.com",
			Password: "",
			FullName: "New User",
		},
		"missing_name": {
			Email:    "newuser@example.com",
			Password: "password123",
			FullName: "",
		},
	}

	loginRequests := map[string]*userentity.LoginRequest{
		"valid": {
			Email:    "user@example.com",
			Password: "password123",
		},
		"missing_email": {
			Email:    "",
			Password: "password123",
		},
		"missing_password": {
			Email:    "user@example.com",
			Password: "",
		},
	}

	createPostRequests := map[string]*postentity.CreatePostRequest{
		"valid_draft": {
			Title:   "Test Draft Post",
			Content: "This is a test draft post content",
			Status:  postentity.PostStatusDraft,
		},
		"valid_published": {
			Title:   "Test Published Post",
			Content: "This is a test published post content",
			Status:  postentity.PostStatusPublished,
		},
		"missing_title": {
			Title:   "",
			Content: "This is a test post content",
			Status:  postentity.PostStatusDraft,
		},
		"missing_content": {
			Title:   "Test Post",
			Content: "",
			Status:  postentity.PostStatusDraft,
		},
		"invalid_status": {
			Title:   "Test Post",
			Content: "This is a test post content",
			Status:  postentity.PostStatus("invalid"),
		},
	}

	updatePostRequests := map[string]*postentity.UpdatePostRequest{
		"update_title": {
			Title:   stringPtr("Updated Title"),
			Content: nil,
			Status:  nil,
		},
		"update_content": {
			Title:   nil,
			Content: stringPtr("Updated content"),
			Status:  nil,
		},
		"update_status": {
			Title:   nil,
			Content: nil,
			Status:  postentityPostStatusPtr(postentity.PostStatusPublished),
		},
		"update_all": {
			Title:   stringPtr("Updated All"),
			Content: stringPtr("Updated all content"),
			Status:  postentityPostStatusPtr(postentity.PostStatusArchived),
		},
	}

	return &RequestFixtures{
		createUserRequests: createUserRequests,
		loginRequests:      loginRequests,
		createPostRequests: createPostRequests,
		updatePostRequests: updatePostRequests,
	}
}

// GetCreateUserRequest returns a create user request fixture by name
func (f *RequestFixtures) GetCreateUserRequest(name string) *userentity.CreateUserRequest {
	if req, exists := f.createUserRequests[name]; exists {
		reqCopy := *req
		return &reqCopy
	}
	return nil
}

// GetLoginRequest returns a login request fixture by name
func (f *RequestFixtures) GetLoginRequest(name string) *userentity.LoginRequest {
	if req, exists := f.loginRequests[name]; exists {
		reqCopy := *req
		return &reqCopy
	}
	return nil
}

// GetCreatePostRequest returns a create post request fixture by name
func (f *RequestFixtures) GetCreatePostRequest(name string) *postentity.CreatePostRequest {
	if req, exists := f.createPostRequests[name]; exists {
		reqCopy := *req
		return &reqCopy
	}
	return nil
}

// GetUpdatePostRequest returns an update post request fixture by name
func (f *RequestFixtures) GetUpdatePostRequest(name string) *postentity.UpdatePostRequest {
	if req, exists := f.updatePostRequests[name]; exists {
		reqCopy := *req
		return &reqCopy
	}
	return nil
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func postentityPostStatusPtr(status postentity.PostStatus) *postentity.PostStatus {
	return &status
}

// TestDatabase provides a clean database state for testing
type TestDatabase struct {
	cleanupFunctions []func()
}

// NewTestDatabase creates a new test database helper
func NewTestDatabase() *TestDatabase {
	return &TestDatabase{
		cleanupFunctions: make([]func(), 0),
	}
}

// AddCleanupFunction adds a function to be called during cleanup
func (td *TestDatabase) AddCleanupFunction(fn func()) {
	td.cleanupFunctions = append(td.cleanupFunctions, fn)
}

// Cleanup executes all cleanup functions
func (td *TestDatabase) Cleanup() {
	// Execute cleanup functions in reverse order
	for i := len(td.cleanupFunctions) - 1; i >= 0; i-- {
		td.cleanupFunctions[i]()
	}
	td.cleanupFunctions = td.cleanupFunctions[:0]
}
