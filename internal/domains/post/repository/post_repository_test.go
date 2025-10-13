package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database/sqlc"
)

// MockPostQuerier is a mock implementation of the Querier interface
type MockPostQuerier struct {
	mock.Mock
}

func (m *MockPostQuerier) CreatePost(ctx context.Context, arg sqlc.CreatePostParams) (sqlc.Post, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.Post), args.Error(1)
}

func (m *MockPostQuerier) GetPostByID(ctx context.Context, id pgtype.UUID) (sqlc.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.Post), args.Error(1)
}

func (m *MockPostQuerier) ListPosts(ctx context.Context, arg sqlc.ListPostsParams) ([]sqlc.Post, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.Post), args.Error(1)
}

func (m *MockPostQuerier) ListPostsByAuthor(ctx context.Context, arg sqlc.ListPostsByAuthorParams) ([]sqlc.Post, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.Post), args.Error(1)
}

func (m *MockPostQuerier) UpdatePost(ctx context.Context, arg sqlc.UpdatePostParams) (sqlc.Post, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.Post), args.Error(1)
}

func (m *MockPostQuerier) DeletePost(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPostRepository_Create(t *testing.T) {
	mockQuerier := new(MockPostQuerier)
	repo := NewPostRepository(mockQuerier)

	ctx := context.Background()
	post := &entity.Post{
		ID:        uuid.New(),
		Title:     "Test Post",
		Content:   "This is a test post content",
		AuthorID:  uuid.New(),
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var idBytes [16]byte
	copy(idBytes[:], post.ID[:])

	var authorIDBytes [16]byte
	copy(authorIDBytes[:], post.AuthorID[:])

	expectedParams := sqlc.CreatePostParams{
		ID:        pgtype.UUID{Bytes: idBytes, Valid: true},
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  pgtype.UUID{Bytes: authorIDBytes, Valid: true},
		Status:    post.Status,
		CreatedAt: pgtype.Timestamptz{Time: post.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}

	expectedDBPost := sqlc.Post{
		ID:        pgtype.UUID{Bytes: idBytes, Valid: true},
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  pgtype.UUID{Bytes: authorIDBytes, Valid: true},
		Status:    post.Status,
		CreatedAt: pgtype.Timestamptz{Time: post.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}

	mockQuerier.On("CreatePost", ctx, expectedParams).Return(expectedDBPost, nil)

	err := repo.Create(ctx, post)
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestPostRepository_GetByID(t *testing.T) {
	mockQuerier := new(MockPostQuerier)
	repo := NewPostRepository(mockQuerier)

	ctx := context.Background()
	postID := uuid.New()
	authorID := uuid.New()

	var idBytes [16]byte
	copy(idBytes[:], postID[:])

	var authorIDBytes [16]byte
	copy(authorIDBytes[:], authorID[:])

	expectedDBPost := sqlc.Post{
		ID:        pgtype.UUID{Bytes: idBytes, Valid: true},
		Title:     "Test Post",
		Content:   "This is test content",
		AuthorID:  pgtype.UUID{Bytes: authorIDBytes, Valid: true},
		Status:    "published",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	mockQuerier.On("GetPostByID", ctx, pgID).Return(expectedDBPost, nil)

	result, err := repo.GetByID(ctx, postID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
	assert.Equal(t, "Test Post", result.Title)
	assert.Equal(t, "This is test content", result.Content)
	assert.Equal(t, authorID, result.AuthorID)
	assert.Equal(t, "published", result.Status)
	mockQuerier.AssertExpectations(t)
}

func TestPostRepository_List(t *testing.T) {
	mockQuerier := new(MockPostQuerier)
	repo := NewPostRepository(mockQuerier)

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedDBPosts := []sqlc.Post{
		{
			ID:        pgtype.UUID{Bytes: [16]byte{1, 2, 3}, Valid: true},
			Title:     "Post 1",
			Content:   "Content 1",
			AuthorID:  pgtype.UUID{Bytes: [16]byte{1, 1, 1}, Valid: true},
			Status:    "published",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        pgtype.UUID{Bytes: [16]byte{4, 5, 6}, Valid: true},
			Title:     "Post 2",
			Content:   "Content 2",
			AuthorID:  pgtype.UUID{Bytes: [16]byte{2, 2, 2}, Valid: true},
			Status:    "draft",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	expectedParams := sqlc.ListPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	mockQuerier.On("ListPosts", ctx, expectedParams).Return(expectedDBPosts, nil)

	result, err := repo.List(ctx, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Post 1", result[0].Title)
	assert.Equal(t, "Post 2", result[1].Title)
	mockQuerier.AssertExpectations(t)
}

func TestPostRepository_ListByAuthor(t *testing.T) {
	mockQuerier := new(MockPostQuerier)
	repo := NewPostRepository(mockQuerier)

	ctx := context.Background()
	authorID := uuid.New()
	limit := 10
	offset := 0

	var authorIDBytes [16]byte
	copy(authorIDBytes[:], authorID[:])

	expectedDBPosts := []sqlc.Post{
		{
			ID:        pgtype.UUID{Bytes: [16]byte{1, 2, 3}, Valid: true},
			Title:     "Author Post 1",
			Content:   "Author Content 1",
			AuthorID:  pgtype.UUID{Bytes: authorIDBytes, Valid: true},
			Status:    "published",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	expectedParams := sqlc.ListPostsByAuthorParams{
		AuthorID: pgtype.UUID{Bytes: authorIDBytes, Valid: true},
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	mockQuerier.On("ListPostsByAuthor", ctx, expectedParams).Return(expectedDBPosts, nil)

	result, err := repo.ListByAuthor(ctx, authorID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Author Post 1", result[0].Title)
	assert.Equal(t, authorID, result[0].AuthorID)
	mockQuerier.AssertExpectations(t)
}

func TestPostRepository_Update(t *testing.T) {
	mockQuerier := new(MockPostQuerier)
	repo := NewPostRepository(mockQuerier)

	ctx := context.Background()
	post := &entity.Post{
		ID:        uuid.New(),
		Title:     "Updated Post",
		Content:   "Updated content",
		Status:    "published",
		UpdatedAt: time.Now(),
	}

	var idBytes [16]byte
	copy(idBytes[:], post.ID[:])

	expectedParams := sqlc.UpdatePostParams{
		ID:        pgtype.UUID{Bytes: idBytes, Valid: true},
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}

	expectedDBPost := sqlc.Post{
		ID:        pgtype.UUID{Bytes: idBytes, Valid: true},
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		UpdatedAt: pgtype.Timestamptz{Time: post.UpdatedAt, Valid: true},
	}

	mockQuerier.On("UpdatePost", ctx, expectedParams).Return(expectedDBPost, nil)

	err := repo.Update(ctx, post)
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestPostRepository_Delete(t *testing.T) {
	mockQuerier := new(MockPostQuerier)
	repo := NewPostRepository(mockQuerier)

	ctx := context.Background()
	postID := uuid.New()

	var idBytes [16]byte
	copy(idBytes[:], postID[:])

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	mockQuerier.On("DeletePost", ctx, pgID).Return(nil)

	err := repo.Delete(ctx, postID)
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}