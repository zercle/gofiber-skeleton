package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skeleton/internal/post/entity"
	mockRepo "github.com/zercle/gofiber-skeleton/internal/post/repository/mocks"
	"github.com/zercle/gofiber-skeleton/internal/post/usecase"
)

func TestPostUsecase_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	threadID := uuid.New()
	userID := uuid.New()
	content := "Test post content"

	expectedPost := &entity.Post{
		ID:        uuid.New(),
		ThreadID:  threadID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockPostRepo.EXPECT().
		Create(ctx, threadID, userID, content).
		Return(expectedPost, nil).
		Times(1)

	result, err := postUsecase.CreatePost(ctx, threadID, userID, content)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, threadID, result.ThreadID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, content, result.Content)
}

func TestPostUsecase_CreatePost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	threadID := uuid.New()
	userID := uuid.New()
	content := "Test post content"

	expectedError := errors.New("database error")

	mockPostRepo.EXPECT().
		Create(ctx, threadID, userID, content).
		Return(nil, expectedError).
		Times(1)

	result, err := postUsecase.CreatePost(ctx, threadID, userID, content)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestPostUsecase_GetPostByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := uuid.New()

	expectedPost := &entity.Post{
		ID:        postID,
		ThreadID:  uuid.New(),
		UserID:    uuid.New(),
		Content:   "Test post content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(expectedPost, nil).
		Times(1)

	result, err := postUsecase.GetPostByID(ctx, postID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
}

func TestPostUsecase_UpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()
	newContent := "Updated content"

	existingPost := &entity.Post{
		ID:        postID,
		ThreadID:  uuid.New(),
		UserID:    userID,
		Content:   "Old content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedPost := &entity.Post{
		ID:        postID,
		ThreadID:  existingPost.ThreadID,
		UserID:    userID,
		Content:   newContent,
		CreatedAt: existingPost.CreatedAt,
		UpdatedAt: time.Now(),
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(existingPost, nil).
		Times(1)

	mockPostRepo.EXPECT().
		Update(ctx, postID, newContent).
		Return(updatedPost, nil).
		Times(1)

	result, err := postUsecase.UpdatePost(ctx, postID, userID, newContent)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newContent, result.Content)
}

func TestPostUsecase_UpdatePost_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := uuid.New()
	actualUserID := uuid.New()
	differentUserID := uuid.New()
	newContent := "Updated content"

	existingPost := &entity.Post{
		ID:        postID,
		ThreadID:  uuid.New(),
		UserID:    actualUserID,
		Content:   "Old content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockPostRepo.EXPECT().
		GetByID(ctx, postID).
		Return(existingPost, nil).
		Times(1)

	result, err := postUsecase.UpdatePost(ctx, postID, differentUserID, newContent)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestPostUsecase_DeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	postID := uuid.New()
	userID := uuid.New()

	existingPost := &entity.Post{
		ID:        postID,
		ThreadID:  uuid.New(),
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

	err := postUsecase.DeletePost(ctx, postID, userID)

	assert.NoError(t, err)
}

func TestPostUsecase_ListPostsByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mockRepo.NewMockPostRepository(ctrl)
	postUsecase := usecase.NewPostUsecase(mockPostRepo)

	ctx := context.Background()
	userID := uuid.New()

	expectedPosts := []*entity.Post{
		{
			ID:        uuid.New(),
			ThreadID:  uuid.New(),
			UserID:    userID,
			Content:   "Post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			ThreadID:  uuid.New(),
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

	result, err := postUsecase.ListPostsByUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, userID, result[0].UserID)
	assert.Equal(t, userID, result[1].UserID)
}
