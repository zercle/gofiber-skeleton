package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domains/post/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/repository"
	"github.com/zercle/gofiber-skeleton/pkg/database"
)

func TestPostRepository_Create_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Mock database expectations
	mock.ExpectExec(`INSERT INTO posts`).
		WithArgs(gomock.Any(), "Test Post", "Test content", "draft", gomock.Any()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Test data
	post := entity.NewPost("Test Post", "Test content", uuid.New())

	// Execute
	err = repo.Create(context.Background(), post)

	// Assert
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetByID_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()
	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "status", "user_id", "created_at", "updated_at"}).
		AddRow(postID, "Test Post", "Test content", "published", userID, time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, title, content, status, user_id, created_at, updated_at FROM posts`).
		WithArgs(postID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetByID(context.Background(), postID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
	assert.Equal(t, "Test Post", result.Title)
	assert.Equal(t, "Test content", result.Content)
	assert.Equal(t, entity.PostStatusPublished, result.Status)
	assert.Equal(t, userID, result.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetByID_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, title, content, status, user_id, created_at, updated_at FROM posts`).
		WithArgs(postID).
		WillReturnError(sql.ErrNoRows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetByID(context.Background(), postID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, entity.ErrPostNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Update_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()
	userID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`UPDATE posts`).
		WithArgs("Updated Title", "Updated Content", "published", gomock.Any(), postID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Test data
	post := &entity.DomainPost{
		ID:        postID,
		Title:     "Updated Title",
		Content:   "Updated Content",
		Status:    entity.PostStatusPublished,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}

	// Execute
	err = repo.Update(context.Background(), post)

	// Assert
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Delete_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()

	// Mock database expectations
	mock.ExpectExec(`DELETE FROM posts`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	err = repo.Delete(context.Background(), postID)

	// Assert
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetByUserID_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "status", "user_id", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Post 1", "Content 1", "draft", userID, time.Now(), time.Now()).
		AddRow(uuid.New(), "Post 2", "Content 2", "published", userID, time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, title, content, status, user_id, created_at, updated_at FROM posts`).
		WithArgs(userID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetByUserID(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Post 1", result[0].Title)
	assert.Equal(t, "Post 2", result[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPublishedPostsByUserID_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "status", "user_id", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Published Post 1", "Content 1", "published", userID, time.Now(), time.Now()).
		AddRow(uuid.New(), "Published Post 2", "Content 2", "published", userID, time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, title, content, status, user_id, created_at, updated_at FROM posts`).
		WithArgs(userID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetPublishedPostsByUserID(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Published Post 1", result[0].Title)
	assert.Equal(t, "Published Post 2", result[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetDraftPostsByUserID_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "status", "user_id", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Draft Post 1", "Content 1", "draft", userID, time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT id, title, content, status, user_id, created_at, updated_at FROM posts`).
		WithArgs(userID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetDraftPostsByUserID(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Draft Post 1", result[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetAllPosts_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "status", "user_id", "full_name", "email", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Post 1", "Content 1", "published", uuid.New(), "Author 1", "author1@example.com", time.Now(), time.Now()).
		AddRow(uuid.New(), "Post 2", "Content 2", "draft", uuid.New(), "Author 2", "author2@example.com", time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT p\.id, p\.title, p\.content, p\.status, p\.user_id, u\.full_name, u\.email, p\.created_at, p\.updated_at`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetAllPosts(context.Background(), 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Post 1", result[0].Title)
	assert.Equal(t, "Author 1", result[0].FullName)
	assert.Equal(t, "Post 2", result[1].Title)
	assert.Equal(t, "Author 2", result[1].FullName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPostsWithAuthor_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "status", "user_id", "full_name", "email", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Published Post", "Content", "published", uuid.New(), "Test Author", "author@example.com", time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT p\.id, p\.title, p\.content, p\.status, p\.user_id, u\.full_name, u\.email, p\.created_at, p\.updated_at`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetPostsWithAuthor(context.Background(), 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Published Post", result[0].Title)
	assert.Equal(t, "Test Author", result[0].FullName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_SearchPosts_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	searchQuery := "test"
	rows := sqlmock.NewRows([]string{"id", "title", "content", "status", "user_id", "full_name", "email", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Test Post", "Test content", "published", uuid.New(), "Test Author", "author@example.com", time.Now(), time.Now())

	// Mock database expectations
	mock.ExpectQuery(`SELECT p\.id, p\.title, p\.content, p\.status, p\.user_id, u\.full_name, u\.email, p\.created_at, p\.updated_at`).
		WithArgs("%"+searchQuery+"%", "published", 10, 0).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.SearchPosts(context.Background(), searchQuery, entity.PostStatusPublished, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test Post", result[0].Title)
	assert.Equal(t, "Test Author", result[0].FullName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetUserPostStats_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := uuid.New()
	lastPostDate := time.Now()
	rows := sqlmock.NewRows([]string{"total_posts", "published_posts", "draft_posts", "last_post_date"}).
		AddRow(10, 7, 3, lastPostDate)

	// Mock database expectations
	mock.ExpectQuery(`SELECT`).
		WithArgs(userID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.GetUserPostStats(context.Background(), userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.TotalPosts)
	assert.Equal(t, 7, result.PublishedPosts)
	assert.Equal(t, 3, result.DraftPosts)
	assert.Equal(t, lastPostDate.Format(time.RFC3339), result.LastPostDate.Format(time.RFC3339))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_CountPostsByStatus_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow("published", 15).
		AddRow("draft", 5).
		AddRow("archived", 2)

	// Mock database expectations
	mock.ExpectQuery(`SELECT status, COUNT\(\*\) as count FROM posts`).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	result, err := repo.CountPostsByStatus(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "published", result[0].Status)
	assert.Equal(t, 15, result[0].Count)
	assert.Equal(t, "draft", result[1].Status)
	assert.Equal(t, 5, result[1].Count)
	assert.Equal(t, "archived", result[2].Status)
	assert.Equal(t, 2, result[2].Count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Exists_True(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()
	rows := sqlmock.NewRows([]string{"exists"}).
		AddRow(true)

	// Mock database expectations
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(postID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	exists, err := repo.Exists(context.Background(), postID)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Exists_False(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()
	rows := sqlmock.NewRows([]string{"exists"}).
		AddRow(false)

	// Mock database expectations
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(postID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	exists, err := repo.Exists(context.Background(), postID)

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_IsOwner_True(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()
	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"is_owner"}).
		AddRow(true)

	// Mock database expectations
	mock.ExpectQuery(`SELECT user_id = \$2 as is_owner FROM posts`).
		WithArgs(postID, userID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	isOwner, err := repo.IsOwner(context.Background(), postID, userID)

	// Assert
	require.NoError(t, err)
	assert.True(t, isOwner)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_IsOwner_False(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postID := uuid.New()
	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"is_owner"}).
		AddRow(false)

	// Mock database expectations
	mock.ExpectQuery(`SELECT user_id = \$2 as is_owner FROM posts`).
		WithArgs(postID, userID).
		WillReturnRows(rows)

	// Create database instance
	sqlxDB := sqlx.NewDb(db, "postgres")
	database := &database.Database{DB: sqlxDB, Pool: &pgxpool.Pool{}}

	// Create test injector
	injector := do.New()
	do.ProvideValue(injector, database)

	// Create repository via constructor
	repo, err := repository.NewPostRepository(injector)
	require.NoError(t, err)

	// Execute
	isOwner, err := repo.IsOwner(context.Background(), postID, userID)

	// Assert
	require.NoError(t, err)
	assert.False(t, isOwner)
	assert.NoError(t, mock.ExpectationsWereMet())
}
