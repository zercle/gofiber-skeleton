package repository

import (
	"context"
	"gofiber-skeleton/internal/entities"
	"gofiber-skeleton/internal/repository/mocks"
	db "gofiber-skeleton/internal/repository/db"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestURLRepository_CreateURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	mockRedis := mocks.NewMockRedisCache(ctrl)
	repo := NewURLRepository(mockQuerier, mockRedis)

	userID := uuid.New()
	url := &entities.URL{
		ID:          uuid.New(),
		OriginalURL: "https://example.com/test",
		ShortCode:   "testcode",
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	mockQuerier.EXPECT().CreateURL(gomock.Any(), db.CreateURLParams{
		OriginalUrl: url.OriginalURL,
		ShortCode:   url.ShortCode,
		UserID:      pgtype.UUID{Bytes: url.UserID, Valid: true},
		ExpiresAt:   pgtype.Timestamptz{Time: url.ExpiresAt, Valid: !url.ExpiresAt.IsZero()},
	}).Return(db.Url{}, nil) // Return empty db.Url and nil error for success

	mockRedis.EXPECT().Set(gomock.Any(), url.ShortCode, url.OriginalURL, time.Duration(0)).Return(redis.NewStatusCmd(context.Background()))

	err := repo.CreateURL(context.Background(), url)
	assert.NoError(t, err)
}

func TestURLRepository_GetURLByShortCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	mockRedis := mocks.NewMockRedisCache(ctrl)
	repo := NewURLRepository(mockQuerier, mockRedis)

	userID := uuid.New()
	expectedURL := &entities.URL{
		ID:          uuid.New(),
		OriginalURL: "https://example.com/get",
		ShortCode:   "getcode",
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	// Test cache hit scenario
	mockRedis.EXPECT().Get(gomock.Any(), expectedURL.ShortCode).Return(redis.NewStringCmd(context.Background())).Times(1).DoAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
		cmd := redis.NewStringCmd(ctx)
		cmd.SetVal(expectedURL.OriginalURL) // Simulate cache hit
		return cmd
	})

	retrievedURLFromCache, err := repo.GetURLByShortCode(context.Background(), expectedURL.ShortCode)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedURLFromCache)
	assert.Equal(t, expectedURL.OriginalURL, retrievedURLFromCache.OriginalURL)
	assert.Equal(t, expectedURL.ShortCode, retrievedURLFromCache.ShortCode) // ShortCode is not stored in Redis cache, so it will be the one passed to GetURLByShortCode

	// Test cache miss scenario
	mockRedis.EXPECT().Get(gomock.Any(), expectedURL.ShortCode).Return(redis.NewStringCmd(context.Background())).Times(1).DoAndReturn(func(ctx context.Context, key string) *redis.StringCmd {
		cmd := redis.NewStringCmd(ctx)
		cmd.SetErr(redis.Nil) // Simulate cache miss
		return cmd
	})
	mockQuerier.EXPECT().GetURLByShortCode(gomock.Any(), expectedURL.ShortCode).Return(db.Url{
		ID:          pgtype.UUID{Bytes: expectedURL.ID, Valid: true},
		OriginalUrl: expectedURL.OriginalURL,
		ShortCode:   expectedURL.ShortCode,
		UserID:      pgtype.UUID{Bytes: expectedURL.UserID, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: expectedURL.CreatedAt, Valid: true},
	}, nil)
	mockRedis.EXPECT().Set(gomock.Any(), expectedURL.ShortCode, expectedURL.OriginalURL, time.Duration(0)).Return(redis.NewStatusCmd(context.Background())).Times(1)

	retrievedURL, err := repo.GetURLByShortCode(context.Background(), expectedURL.ShortCode)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedURL)
	assert.Equal(t, expectedURL.OriginalURL, retrievedURL.OriginalURL)
	assert.Equal(t, expectedURL.ShortCode, retrievedURL.ShortCode)
	assert.Equal(t, expectedURL.ID, retrievedURL.ID)
	assert.Equal(t, expectedURL.UserID, retrievedURL.UserID)
}

func TestURLRepository_GetURLsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	mockRedis := mocks.NewMockRedisCache(ctrl)
	repo := NewURLRepository(mockQuerier, mockRedis)

	userID := uuid.New()
	expectedURLs := []db.Url{
		{
			ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
			OriginalUrl: "https://example.com/user1",
			ShortCode:   "user1",
			UserID:      pgtype.UUID{Bytes: userID, Valid: true},
			CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
			OriginalUrl: "https://example.com/user2",
			ShortCode:   "user2",
			UserID:      pgtype.UUID{Bytes: userID, Valid: true},
			CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	mockQuerier.EXPECT().GetURLsByUserID(gomock.Any(), pgtype.UUID{Bytes: userID, Valid: true}).Return(expectedURLs, nil)

	retrievedURLs, err := repo.GetURLsByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, retrievedURLs, 2)
	assert.Equal(t, expectedURLs[0].OriginalUrl, retrievedURLs[0].OriginalURL)
	assert.Equal(t, expectedURLs[1].OriginalUrl, retrievedURLs[1].OriginalURL)
}

func TestURLRepository_UpdateURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	mockRedis := mocks.NewMockRedisCache(ctrl)
	repo := NewURLRepository(mockQuerier, mockRedis)

	urlID := uuid.New()
	newOriginalURL := "https://example.com/updated"

	mockQuerier.EXPECT().UpdateURL(gomock.Any(), db.UpdateURLParams{
		ID:          pgtype.UUID{Bytes: urlID, Valid: true},
		OriginalUrl: newOriginalURL,
	}).Return(db.Url{}, nil)

	err := repo.UpdateURL(context.Background(), &entities.URL{ID: urlID, OriginalURL: newOriginalURL})
	assert.NoError(t, err)
}

func TestURLRepository_DeleteURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	mockRedis := mocks.NewMockRedisCache(ctrl)
	repo := NewURLRepository(mockQuerier, mockRedis)

	urlID := uuid.New()

	mockQuerier.EXPECT().DeleteURL(gomock.Any(), pgtype.UUID{Bytes: urlID, Valid: true}).Return(nil)

	err := repo.DeleteURL(context.Background(), urlID)
	assert.NoError(t, err)
}

func TestURLRepository_GetURLByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	mockRedis := mocks.NewMockRedisCache(ctrl)
	repo := NewURLRepository(mockQuerier, mockRedis)

	userID := uuid.New()
	expectedURL := &entities.URL{
		ID:          uuid.New(),
		OriginalURL: "https://example.com/getbyid",
		ShortCode:   "getbyid",
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	mockQuerier.EXPECT().GetURLByID(gomock.Any(), pgtype.UUID{Bytes: expectedURL.ID, Valid: true}).Return(db.Url{
		ID:          pgtype.UUID{Bytes: expectedURL.ID, Valid: true},
		OriginalUrl: expectedURL.OriginalURL,
		ShortCode:   expectedURL.ShortCode,
		UserID:      pgtype.UUID{Bytes: expectedURL.UserID, Valid: true},
		CreatedAt:   pgtype.Timestamptz{Time: expectedURL.CreatedAt, Valid: true},
	}, nil)

	retrievedURL, err := repo.GetURLByID(context.Background(), expectedURL.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedURL)
	assert.Equal(t, expectedURL.OriginalURL, retrievedURL.OriginalURL)
	assert.Equal(t, expectedURL.ShortCode, retrievedURL.ShortCode)
	assert.Equal(t, expectedURL.ID, retrievedURL.ID)
	assert.Equal(t, expectedURL.UserID, retrievedURL.UserID)
}