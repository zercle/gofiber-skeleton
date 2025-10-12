package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/mocks"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/usecase"
)

func TestPostUsecase_CreatePost_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	req := &entity.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
		Status:  entity.PostStatusDraft,
	}

	// Setup expectations
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	result, err := uc.CreatePost(context.Background(), userID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)
	assert.Equal(t, req.Status, result.Status)
	assert.Equal(t, userID, result.UserID)
}

func TestPostUsecase_CreatePost_MissingTitle(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	req := &entity.CreatePostRequest{
		Title:   "",
		Content: "This is a test post content",
		Status:  entity.PostStatusDraft,
	}

	// Execute
	result, err := uc.CreatePost(context.Background(), userID, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrPostTitleRequired, err)
}

func TestPostUsecase_CreatePost_MissingContent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	req := &entity.CreatePostRequest{
		Title:   "Test Post",
		Content: "",
		Status:  entity.PostStatusDraft,
	}

	// Execute
	result, err := uc.CreatePost(context.Background(), userID, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrPostContentRequired, err)
}

func TestPostUsecase_CreatePost_InvalidStatus(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	req := &entity.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
		Status:  entity.PostStatus("invalid"),
	}

	// Execute
	result, err := uc.CreatePost(context.Background(), userID, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrInvalidPostStatus, err)
}

func TestPostUsecase_GetPost_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()
	expectedPost := &entity.DomainPost{
		ID:        postID,
		Title:     "Test Post",
		Content:   "This is a test post content",
		Status:    entity.PostStatusPublished,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(expectedPost, nil)

	// Execute
	result, err := uc.GetPost(context.Background(), postID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPost.ID, result.ID)
	assert.Equal(t, expectedPost.Title, result.Title)
	assert.Equal(t, expectedPost.Content, result.Content)
}

func TestPostUsecase_GetPost_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()

	// Setup expectations
	mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(nil, entity.ErrPostNotFound)

	// Execute
	result, err := uc.GetPost(context.Background(), postID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrPostNotFound, err)
}

func TestPostUsecase_UpdatePost_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()
	existingPost := &entity.DomainPost{
		ID:        postID,
		Title:     "Original Title",
		Content:   "Original content",
		Status:    entity.PostStatusDraft,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &entity.UpdatePostRequest{
		Title:   stringPtr("Updated Title"),
		Content: stringPtr("Updated content"),
		Status:  entityPostStatusPtr(entity.PostStatusPublished),
	}

	// Setup expectations
	mockRepo.EXPECT().IsOwner(gomock.Any(), postID, userID).Return(true, nil)
	mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	result, err := uc.UpdatePost(context.Background(), postID, userID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, *req.Title, result.Title)
	assert.Equal(t, *req.Content, result.Content)
	assert.Equal(t, *req.Status, result.Status)
}

func TestPostUsecase_UpdatePost_NotOwner(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()
	req := &entity.UpdatePostRequest{
		Title:   stringPtr("Updated Title"),
		Content: stringPtr("Updated content"),
	}

	// Setup expectations
	mockRepo.EXPECT().IsOwner(gomock.Any(), postID, userID).Return(false, nil)

	// Execute
	result, err := uc.UpdatePost(context.Background(), postID, userID, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrPostNotFound, err)
}

func TestPostUsecase_DeletePost_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()

	// Setup expectations
	mockRepo.EXPECT().IsOwner(gomock.Any(), postID, userID).Return(true, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), postID).Return(nil)

	// Execute
	err = uc.DeletePost(context.Background(), postID, userID)

	// Assert
	require.NoError(t, err)
}

func TestPostUsecase_DeletePost_NotOwner(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()

	// Setup expectations
	mockRepo.EXPECT().IsOwner(gomock.Any(), postID, userID).Return(false, nil)

	// Execute
	err = uc.DeletePost(context.Background(), postID, userID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, entity.ErrPostNotFound, err)
}

func TestPostUsecase_GetUserPosts_WithStatus(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	status := entity.PostStatusPublished
	expectedPosts := []*entity.DomainPost{
		{
			ID:        uuid.New(),
			Title:     "Published Post 1",
			Content:   "Content 1",
			Status:    entity.PostStatusPublished,
			UserID:    userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Published Post 2",
			Content:   "Content 2",
			Status:    entity.PostStatusPublished,
			UserID:    userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Setup expectations
	mockRepo.EXPECT().GetPublishedPostsByUserID(gomock.Any(), userID).Return(expectedPosts, nil)

	// Execute
	result, err := uc.GetUserPosts(context.Background(), userID, &status)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedPosts[0].Title, result[0].Title)
	assert.Equal(t, expectedPosts[1].Title, result[1].Title)
}

func TestPostUsecase_GetUserPosts_AllPosts(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	expectedPosts := []*entity.DomainPost{
		{
			ID:        uuid.New(),
			Title:     "Post 1",
			Content:   "Content 1",
			Status:    entity.PostStatusDraft,
			UserID:    userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Setup expectations
	mockRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(expectedPosts, nil)

	// Execute
	result, err := uc.GetUserPosts(context.Background(), userID, nil)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedPosts[0].Title, result[0].Title)
}

func TestPostUsecase_GetUserPostStats_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	userID := uuid.New()
	expectedStats := &entity.PostStats{
		TotalPosts:     10,
		PublishedPosts: 7,
		DraftPosts:     3,
		LastPostDate:   &time.Time{},
	}

	// Setup expectations
	mockRepo.EXPECT().GetUserPostStats(gomock.Any(), userID).Return(expectedStats, nil)

	// Execute
	result, err := uc.GetUserPostStats(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStats.TotalPosts, result.TotalPosts)
	assert.Equal(t, expectedStats.PublishedPosts, result.PublishedPosts)
	assert.Equal(t, expectedStats.DraftPosts, result.DraftPosts)
}

func TestPostUsecase_GetAllPosts_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	expectedPosts := []*entity.PostWithAuthor{
		{
			ID:        uuid.New(),
			Title:     "Post 1",
			Content:   "Content 1",
			Status:    entity.PostStatusPublished,
			UserID:    uuid.New(),
			FullName:  "Author 1",
			Email:     "author1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Setup expectations
	mockRepo.EXPECT().GetAllPosts(gomock.Any(), 10, 0).Return(expectedPosts, nil)

	// Execute
	result, err := uc.GetAllPosts(context.Background(), 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedPosts[0].Title, result[0].Title)
	assert.Equal(t, expectedPosts[0].FullName, result[0].FullName)
}

func TestPostUsecase_SearchPosts_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	searchQuery := "test"
	expectedPosts := []*entity.PostWithAuthor{
		{
			ID:        uuid.New(),
			Title:     "Test Post",
			Content:   "This is a test post",
			Status:    entity.PostStatusPublished,
			UserID:    uuid.New(),
			FullName:  "Test Author",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Setup expectations
	mockRepo.EXPECT().SearchPosts(gomock.Any(), searchQuery, entity.PostStatusPublished, 10, 0).Return(expectedPosts, nil)

	// Execute
	result, err := uc.SearchPosts(context.Background(), searchQuery, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedPosts[0].Title, result[0].Title)
}

func TestPostUsecase_SearchPosts_EmptyQuery(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Execute
	result, err := uc.SearchPosts(context.Background(), "", 10, 0)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "search query is required")
}

func TestPostUsecase_PublishPost_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()
	existingPost := &entity.DomainPost{
		ID:        postID,
		Title:     "Test Post",
		Content:   "Test content",
		Status:    entity.PostStatusDraft,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockRepo.EXPECT().IsOwner(gomock.Any(), postID, userID).Return(true, nil)
	mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	err = uc.PublishPost(context.Background(), postID, userID)

	// Assert
	require.NoError(t, err)
}

func TestPostUsecase_ArchivePost_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()
	existingPost := &entity.DomainPost{
		ID:        postID,
		Title:     "Test Post",
		Content:   "Test content",
		Status:    entity.PostStatusPublished,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockRepo.EXPECT().IsOwner(gomock.Any(), postID, userID).Return(true, nil)
	mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	err = uc.ArchivePost(context.Background(), postID, userID)

	// Assert
	require.NoError(t, err)
}

func TestPostUsecase_UnpublishPost_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	// Create a test usecase instance using dependency injection pattern
	injector := do.New()
	do.ProvideValue(injector, mockRepo)

	uc, err := usecase.NewPostUsecase(injector)
	require.NoError(t, err)

	// Test data
	postID := uuid.New()
	userID := uuid.New()
	existingPost := &entity.DomainPost{
		ID:        postID,
		Title:     "Test Post",
		Content:   "Test content",
		Status:    entity.PostStatusPublished,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup expectations
	mockRepo.EXPECT().IsOwner(gomock.Any(), postID, userID).Return(true, nil)
	mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	// Execute
	err = uc.UnpublishPost(context.Background(), postID, userID)

	// Assert
	require.NoError(t, err)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func entityPostStatusPtr(status entity.PostStatus) *entity.PostStatus {
	return &status
}
