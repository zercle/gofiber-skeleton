package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database/sqlc"
)

// MockQuerier is a mock implementation of the Querier interface
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) GetUserByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockQuerier) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) ListUsers(ctx context.Context, arg sqlc.ListUsersParams) ([]sqlc.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.User), args.Error(1)
}

func TestUserRepository_Create(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := NewUserRepository(mockQuerier)

	ctx := context.Background()
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	var idBytes [16]byte
	copy(idBytes[:], user.ID[:])

	expectedParams := sqlc.CreateUserParams{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		CreatedAt:    pgtype.Timestamptz{Time: user.CreatedAt, Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: user.UpdatedAt, Valid: true},
	}

	expectedDBUser := sqlc.User{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		CreatedAt:    pgtype.Timestamptz{Time: user.CreatedAt, Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: user.UpdatedAt, Valid: true},
	}

	mockQuerier.On("CreateUser", ctx, expectedParams).Return(expectedDBUser, nil)

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestUserRepository_GetByID(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := NewUserRepository(mockQuerier)

	ctx := context.Background()
	userID := uuid.New()

	var idBytes [16]byte
	copy(idBytes[:], userID[:])

	expectedDBUser := sqlc.User{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	mockQuerier.On("GetUserByID", ctx, pgID).Return(expectedDBUser, nil)

	result, err := repo.GetByID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.FullName)
	mockQuerier.AssertExpectations(t)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := NewUserRepository(mockQuerier)

	ctx := context.Background()
	email := "test@example.com"
	userID := uuid.New()

	var idBytes [16]byte
	copy(idBytes[:], userID[:])

	expectedDBUser := sqlc.User{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        email,
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	mockQuerier.On("GetUserByEmail", ctx, email).Return(expectedDBUser, nil)

	result, err := repo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, "Test User", result.FullName)
	mockQuerier.AssertExpectations(t)
}

func TestUserRepository_List(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := NewUserRepository(mockQuerier)

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedDBUsers := []sqlc.User{
		{
			ID:           pgtype.UUID{Bytes: [16]byte{1, 2, 3}, Valid: true},
			Email:        "user1@example.com",
			PasswordHash: "hash1",
			FullName:     "User One",
			CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:           pgtype.UUID{Bytes: [16]byte{4, 5, 6}, Valid: true},
			Email:        "user2@example.com",
			PasswordHash: "hash2",
			FullName:     "User Two",
			CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	expectedParams := sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	mockQuerier.On("ListUsers", ctx, expectedParams).Return(expectedDBUsers, nil)

	result, err := repo.List(ctx, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "user1@example.com", result[0].Email)
	assert.Equal(t, "user2@example.com", result[1].Email)
	mockQuerier.AssertExpectations(t)
}

func TestUserRepository_Update(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := NewUserRepository(mockQuerier)

	ctx := context.Background()
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "updated@example.com",
		PasswordHash: "updatedhash",
		FullName:     "Updated User",
		UpdatedAt:    time.Now(),
	}

	var idBytes [16]byte
	copy(idBytes[:], user.ID[:])

	expectedParams := sqlc.UpdateUserParams{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		UpdatedAt:    pgtype.Timestamptz{Time: user.UpdatedAt, Valid: true},
	}

	expectedDBUser := sqlc.User{
		ID:           pgtype.UUID{Bytes: idBytes, Valid: true},
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		UpdatedAt:    pgtype.Timestamptz{Time: user.UpdatedAt, Valid: true},
	}

	mockQuerier.On("UpdateUser", ctx, expectedParams).Return(expectedDBUser, nil)

	err := repo.Update(ctx, user)
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestUserRepository_Delete(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := NewUserRepository(mockQuerier)

	ctx := context.Background()
	userID := uuid.New()

	var idBytes [16]byte
	copy(idBytes[:], userID[:])

	pgID := pgtype.UUID{Bytes: idBytes, Valid: true}
	mockQuerier.On("DeleteUser", ctx, pgID).Return(nil)

	err := repo.Delete(ctx, userID)
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}