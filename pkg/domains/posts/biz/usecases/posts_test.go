package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/pkg/domains/posts/entities"
	"github.com/zercle/gofiber-skeleton/pkg/domains/posts/mocks"
	"github.com/zercle/gofiber-skeleton/pkg/domains/posts/models"
)

func TestPostUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	userID := uuid.New()
	req := &models.CreatePostRequest{
		Title:       "Test Post",
		Content:     "This is test content for the post",
		IsPublished: true,
	}

	t.Run("Successful post creation", func(t *testing.T) {
		expectedPost := &entities.Post{
			ID:          uuid.New(),
			UserID:      userID,
			Title:       req.Title,
			Content:     req.Content,
			Slug:        "test-post-1234567890",
			IsPublished: req.IsPublished,
			PublishedAt: &time.Time{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedPost, nil).Times(1)

		resp, err := uc.Create(context.Background(), userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedPost.ID, resp.ID)
		assert.Equal(t, expectedPost.Title, resp.Title)
		assert.Equal(t, expectedPost.UserID, resp.UserID)
		assert.Equal(t, req.IsPublished, resp.IsPublished)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)

		resp, err := uc.Create(context.Background(), userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "failed to create post")
	})
}

func TestPostUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	postID := uuid.New()

	t.Run("Successful post retrieval", func(t *testing.T) {
		expectedPost := &entities.Post{
			ID:          postID,
			UserID:      uuid.New(),
			Title:       "Test Post",
			Content:     "Test content",
			Slug:        "test-post",
			IsPublished: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(expectedPost, nil).Times(1)

		resp, err := uc.GetByID(context.Background(), postID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedPost.ID, resp.ID)
		assert.Equal(t, expectedPost.Title, resp.Title)
	})

	t.Run("Post not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(nil, errors.New("not found")).Times(1)

		resp, err := uc.GetByID(context.Background(), postID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "post not found")
	})
}

func TestPostUseCase_GetBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	slug := "test-post-slug"

	t.Run("Successful post retrieval by slug", func(t *testing.T) {
		expectedPost := &entities.Post{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Title:       "Test Post",
			Content:     "Test content",
			Slug:        slug,
			IsPublished: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetBySlug(gomock.Any(), slug).Return(expectedPost, nil).Times(1)

		resp, err := uc.GetBySlug(context.Background(), slug)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedPost.Slug, resp.Slug)
		assert.Equal(t, expectedPost.Title, resp.Title)
	})

	t.Run("Post not found by slug", func(t *testing.T) {
		mockRepo.EXPECT().GetBySlug(gomock.Any(), slug).Return(nil, errors.New("not found")).Times(1)

		resp, err := uc.GetBySlug(context.Background(), slug)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "post not found")
	})
}

func TestPostUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	limit := 10
	offset := 0
	isPublished := true

	t.Run("Successful posts listing", func(t *testing.T) {
		posts := []*entities.Post{
			{
				ID:          uuid.New(),
				UserID:      uuid.New(),
				Title:       "Post 1",
				Content:     "Content 1",
				Slug:        "post-1",
				IsPublished: true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          uuid.New(),
				UserID:      uuid.New(),
				Title:       "Post 2",
				Content:     "Content 2",
				Slug:        "post-2",
				IsPublished: true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockRepo.EXPECT().List(gomock.Any(), limit, offset, &isPublished).Return(posts, nil).Times(1)

		resp, err := uc.List(context.Background(), limit, offset, &isPublished)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Posts, 2)
		assert.Equal(t, limit, resp.Limit)
		assert.Equal(t, offset, resp.Offset)
		assert.Equal(t, 2, resp.Total)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.EXPECT().List(gomock.Any(), limit, offset, &isPublished).Return(nil, errors.New("db error")).Times(1)

		resp, err := uc.List(context.Background(), limit, offset, &isPublished)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "failed to fetch posts")
	})
}

func TestPostUseCase_ListByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	userID := uuid.New()
	limit := 10
	offset := 0

	t.Run("Successful user posts listing", func(t *testing.T) {
		posts := []*entities.Post{
			{
				ID:          uuid.New(),
				UserID:      userID,
				Title:       "User Post 1",
				Content:     "User Content 1",
				Slug:        "user-post-1",
				IsPublished: true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockRepo.EXPECT().ListByUser(gomock.Any(), userID, limit, offset).Return(posts, nil).Times(1)

		resp, err := uc.ListByUser(context.Background(), userID, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Posts, 1)
		assert.Equal(t, userID, resp.Posts[0].UserID)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.EXPECT().ListByUser(gomock.Any(), userID, limit, offset).Return(nil, errors.New("db error")).Times(1)

		resp, err := uc.ListByUser(context.Background(), userID, limit, offset)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "failed to fetch user posts")
	})
}

func TestPostUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	postID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	existingPost := &entities.Post{
		ID:          postID,
		UserID:      userID,
		Title:       "Original Title",
		Content:     "Original Content",
		Slug:        "original-title",
		IsPublished: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("Successful post update", func(t *testing.T) {
		newTitle := "Updated Title"
		newContent := "Updated Content"
		isPublished := true

		req := &models.UpdatePostRequest{
			Title:       &newTitle,
			Content:     &newContent,
			IsPublished: &isPublished,
		}

		updatedPost := &entities.Post{
			ID:          postID,
			UserID:      userID,
			Title:       newTitle,
			Content:     newContent,
			Slug:        "updated-title-1234567890",
			IsPublished: isPublished,
			PublishedAt: &time.Time{},
			CreatedAt:   existingPost.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updatedPost, nil).Times(1)

		resp, err := uc.Update(context.Background(), postID, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, newTitle, resp.Title)
		assert.Equal(t, newContent, resp.Content)
		assert.Equal(t, isPublished, resp.IsPublished)
	})

	t.Run("Post not found", func(t *testing.T) {
		req := &models.UpdatePostRequest{}

		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(nil, errors.New("not found")).Times(1)

		resp, err := uc.Update(context.Background(), postID, userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "post not found")
	})

	t.Run("Unauthorized update", func(t *testing.T) {
		req := &models.UpdatePostRequest{}

		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)

		resp, err := uc.Update(context.Background(), postID, otherUserID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "unauthorized: you can only update your own posts")
	})

	t.Run("Repository update error", func(t *testing.T) {
		req := &models.UpdatePostRequest{}

		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error")).Times(1)

		resp, err := uc.Update(context.Background(), postID, userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "failed to update post")
	})
}

func TestPostUseCase_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	postID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	existingPost := &entities.Post{
		ID:          postID,
		UserID:      userID,
		Title:       "Test Post",
		Content:     "Test Content",
		Slug:        "test-post",
		IsPublished: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("Successful post publish", func(t *testing.T) {
		publishedPost := &entities.Post{
			ID:          postID,
			UserID:      userID,
			Title:       "Test Post",
			Content:     "Test Content",
			Slug:        "test-post",
			IsPublished: true,
			PublishedAt: &time.Time{},
			CreatedAt:   existingPost.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Publish(gomock.Any(), postID).Return(publishedPost, nil).Times(1)

		resp, err := uc.Publish(context.Background(), postID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.IsPublished)
	})

	t.Run("Post not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(nil, errors.New("not found")).Times(1)

		resp, err := uc.Publish(context.Background(), postID, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "post not found")
	})

	t.Run("Unauthorized publish", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)

		resp, err := uc.Publish(context.Background(), postID, otherUserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "unauthorized: you can only publish your own posts")
	})

	t.Run("Repository publish error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Publish(gomock.Any(), postID).Return(nil, errors.New("db error")).Times(1)

		resp, err := uc.Publish(context.Background(), postID, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "failed to publish post")
	})
}

func TestPostUseCase_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	postID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	existingPost := &entities.Post{
		ID:          postID,
		UserID:      userID,
		Title:       "Test Post",
		Content:     "Test Content",
		Slug:        "test-post",
		IsPublished: true,
		PublishedAt: &time.Time{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("Successful post unpublish", func(t *testing.T) {
		unpublishedPost := &entities.Post{
			ID:          postID,
			UserID:      userID,
			Title:       "Test Post",
			Content:     "Test Content",
			Slug:        "test-post",
			IsPublished: false,
			PublishedAt: nil,
			CreatedAt:   existingPost.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Unpublish(gomock.Any(), postID).Return(unpublishedPost, nil).Times(1)

		resp, err := uc.Unpublish(context.Background(), postID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.IsPublished)
	})

	t.Run("Post not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(nil, errors.New("not found")).Times(1)

		resp, err := uc.Unpublish(context.Background(), postID, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "post not found")
	})

	t.Run("Unauthorized unpublish", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)

		resp, err := uc.Unpublish(context.Background(), postID, otherUserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "unauthorized: you can only unpublish your own posts")
	})

	t.Run("Repository unpublish error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Unpublish(gomock.Any(), postID).Return(nil, errors.New("db error")).Times(1)

		resp, err := uc.Unpublish(context.Background(), postID, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "failed to unpublish post")
	})
}

func TestPostUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo)

	postID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	existingPost := &entities.Post{
		ID:          postID,
		UserID:      userID,
		Title:       "Test Post",
		Content:     "Test Content",
		Slug:        "test-post",
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("Successful post deletion", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Delete(gomock.Any(), postID).Return(nil).Times(1)

		err := uc.Delete(context.Background(), postID, userID)

		assert.NoError(t, err)
	})

	t.Run("Post not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(nil, errors.New("not found")).Times(1)

		err := uc.Delete(context.Background(), postID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "post not found")
	})

	t.Run("Unauthorized delete", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)

		err := uc.Delete(context.Background(), postID, otherUserID)

		assert.Error(t, err)
		assert.EqualError(t, err, "unauthorized: you can only delete your own posts")
	})

	t.Run("Repository delete error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), postID).Return(existingPost, nil).Times(1)
		mockRepo.EXPECT().Delete(gomock.Any(), postID).Return(errors.New("db error")).Times(1)

		err := uc.Delete(context.Background(), postID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to delete post")
	})
}

func TestPostUseCase_generateSlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	uc := NewPostUseCase(mockRepo).(*postUseCase)

	t.Run("Generate slug from title", func(t *testing.T) {
		title := "This Is A Test Title"
		slug := uc.generateSlug(title)

		assert.Contains(t, slug, "this-is-a-test-title")
		assert.Contains(t, slug, "-") // Should contain timestamp
	})

	t.Run("Generate slug with special characters", func(t *testing.T) {
		title := "Test Title With Special Characters!"
		slug := uc.generateSlug(title)

		assert.Contains(t, slug, "test-title-with-special-characters!")
	})
}