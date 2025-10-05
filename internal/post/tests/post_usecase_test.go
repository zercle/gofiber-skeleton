package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/post/entity"
	mockRepo "github.com/zercle/gofiber-skeleton/internal/post/repository/mocks"
	"github.com/zercle/gofiber-skeleton/internal/post/usecase"
)

// mustNewUUID is a test helper that generates a UUIDv7 or fails the test
func mustNewUUID(t *testing.T) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV7()
	require.NoError(t, err)
	return id
}

func TestPostUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	threadID := mustNewUUID(t)
	userID := mustNewUUID(t)
	content := "Test post content"

	mockPostRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, p *entity.Post) error {
			// Simulate database setting the ID
			p.ID = mustNewUUID(t)
			p.CreatedAt = time.Now()
			p.UpdatedAt = time.Now()
			return nil
		}).
		Times(1)

	result, err := postUsecase.Create(ctx, userID, threadID, content)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, threadID, result.ThreadID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, content, result.Content)
	assert.NotEqual(t, uuid.Nil, result.ID)
}

func TestPostUsecase_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	threadID := mustNewUUID(t)
	userID := mustNewUUID(t)
	content := "Test post content"

	expectedError := errors.New("database error")

	mockPostRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(expectedError).
		Times(1)

	result, err := postUsecase.Create(ctx, userID, threadID, content)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create post")
}

func TestPostUsecase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := mustNewUUID(t)

	expectedPost := &entity.Post{
		ID:        postID,
		ThreadID:  mustNewUUID(t),
		UserID:    mustNewUUID(t),
		Content:   "Test post content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(expectedPost, nil).
		Times(1)

	result, err := postUsecase.Get(ctx, postID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
}

func TestPostUsecase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := mustNewUUID(t)
	userID := mustNewUUID(t)
	newContent := "Updated content"

	existingPost := &entity.Post{
		ID:        postID,
		ThreadID:  mustNewUUID(t),
		UserID:    userID,
		Content:   "Old content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatePost := &entity.Post{
		ID:      postID,
		Content: newContent,
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(existingPost, nil).
		Times(1)

	mockPostRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, p *entity.Post) error {
			assert.Equal(t, newContent, p.Content)
			assert.Equal(t, userID, p.UserID)
			return nil
		}).
		Times(1)

	result, err := postUsecase.Update(ctx, userID, updatePost)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newContent, result.Content)
}

func TestPostUsecase_Update_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := mustNewUUID(t)
	actualUserID := mustNewUUID(t)
	differentUserID := mustNewUUID(t)
	newContent := "Updated content"

	existingPost := &entity.Post{
		ID:        postID,
		ThreadID:  mustNewUUID(t),
		UserID:    actualUserID,
		Content:   "Old content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatePost := &entity.Post{
		ID:      postID,
		Content: newContent,
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(existingPost, nil).
		Times(1)

	result, err := postUsecase.Update(ctx, differentUserID, updatePost)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not the owner")
}

func TestPostUsecase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := mustNewUUID(t)
	userID := mustNewUUID(t)

	existingPost := &entity.Post{
		ID:        postID,
		ThreadID:  mustNewUUID(t),
		UserID:    userID,
		Content:   "Test content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(existingPost, nil).
		Times(1)

	mockPostRepo.EXPECT().
		Delete(ctx, postID).
		Return(nil).
		Times(1)

	err := postUsecase.Delete(ctx, userID, postID)

	assert.NoError(t, err)
}

func TestPostUsecase_Delete_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := mustNewUUID(t)
	actualUserID := mustNewUUID(t)
	differentUserID := mustNewUUID(t)

	existingPost := &entity.Post{
		ID:        postID,
		ThreadID:  mustNewUUID(t),
		UserID:    actualUserID,
		Content:   "Test content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(existingPost, nil).
		Times(1)

	err := postUsecase.Delete(ctx, differentUserID, postID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not the owner")
}

func TestPostUsecase_ListByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	userID := mustNewUUID(t)

	expectedPosts := []*entity.Post{
		{
			ID:        mustNewUUID(t),
			ThreadID:  mustNewUUID(t),
			UserID:    userID,
			Content:   "Post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        mustNewUUID(t),
			ThreadID:  mustNewUUID(t),
			UserID:    userID,
			Content:   "Post 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockPostRepo.EXPECT().
		ListByUser(ctx, userID).
		Return(expectedPosts, nil).
		Times(1)

	result, err := postUsecase.ListByUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, userID, result[0].UserID)
	assert.Equal(t, userID, result[1].UserID)
}
