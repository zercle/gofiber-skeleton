package repository_test

import (
	"context"
	"gofiber-skeleton/internal/entities"
	repo "gofiber-skeleton/internal/repository"
	"gofiber-skeleton/internal/repository/mocks"
	db "gofiber-skeleton/internal/repository/db"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserRepository_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	repo := repo.NewSQLUserRepository(mockQuerier)

	user := &entities.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockQuerier.EXPECT().CreateUser(gomock.Any(), db.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
	}).Return(db.User{}, nil)

	err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	repo := repo.NewSQLUserRepository(mockQuerier)

	expectedUser := &entities.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockQuerier.EXPECT().GetUserByID(gomock.Any(), pgtype.UUID{Bytes: expectedUser.ID, Valid: true}).Return(db.User{
		ID:        pgtype.UUID{Bytes: expectedUser.ID, Valid: true},
		Username:  expectedUser.Username,
		Password:  expectedUser.Password,
		CreatedAt: pgtype.Timestamptz{Time: expectedUser.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: expectedUser.UpdatedAt, Valid: true},
	}, nil)

	retrievedUser, err := repo.GetUserByID(context.Background(), expectedUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, expectedUser.Username, retrievedUser.Username)
	assert.Equal(t, expectedUser.ID, retrievedUser.ID)
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mocks.NewMockDBQueriesInterface(ctrl)
	repo := repo.NewSQLUserRepository(mockQuerier)

	expectedUser := &entities.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockQuerier.EXPECT().GetUserByUsername(gomock.Any(), expectedUser.Username).Return(db.User{
		ID:        pgtype.UUID{Bytes: expectedUser.ID, Valid: true},
		Username:  expectedUser.Username,
		Password:  expectedUser.Password,
		CreatedAt: pgtype.Timestamptz{Time: expectedUser.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: expectedUser.UpdatedAt, Valid: true},
	}, nil)

	retrievedUser, err := repo.GetUserByUsername(context.Background(), expectedUser.Username)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, expectedUser.Username, retrievedUser.Username)
	assert.Equal(t, expectedUser.ID, retrievedUser.ID)
}