package auth_repositories

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/entities"
	sqldb "github.com/zercle/gofiber-skeleton/internal/infrastructure/database/queries"
	mock_database "github.com/zercle/gofiber-skeleton/internal/infrastructure/database/queries/mocks"
)

func TestUserRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_database.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()
	testUser := &entities.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "John",
		LastName:     "Doe",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sqlcUser := sqldb.Users{
		ID:              pgtype.UUID{Bytes: userID, Valid: true},
		Email:           testUser.Email,
		PasswordHash:    testUser.PasswordHash,
		FirstName:       testUser.FirstName,
		LastName:        testUser.LastName,
		IsActive:        pgtype.Bool{Bool: testUser.IsActive, Valid: true},
		IsEmailVerified: pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: testUser.CreatedAt, Valid: true},
		UpdatedAt:       pgtype.Timestamptz{Time: testUser.UpdatedAt, Valid: true},
	}

	createParams := sqldb.CreateUserParams{
		Email:        testUser.Email,
		PasswordHash: testUser.PasswordHash,
		FirstName:    testUser.FirstName,
		LastName:     testUser.LastName,
	}

	mockQuerier.EXPECT().CreateUser(ctx, createParams).Return(sqlcUser, nil)

	createdUser, err := repo.Create(ctx, testUser)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, testUser.Email, createdUser.Email)
	assert.Equal(t, testUser.FirstName, createdUser.FirstName)
}

func TestUserRepository_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()
	testUser := &entities.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "John",
		LastName:     "Doe",
	}

	createParams := sqldb.CreateUserParams{
		Email:        testUser.Email,
		PasswordHash: testUser.PasswordHash,
		FirstName:    testUser.FirstName,
		LastName:     testUser.LastName,
	}

	mockQuerier.EXPECT().CreateUser(ctx, createParams).Return(sqldb.Users{}, fmt.Errorf("db error"))

	createdUser, err := repo.Create(ctx, testUser)
	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestUserRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()
	testUser := &entities.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "John",
		LastName:     "Doe",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sqlcUser := sqldb.Users{
		ID:              pgtype.UUID{Bytes: userID, Valid: true},
		Email:           testUser.Email,
		PasswordHash:    testUser.PasswordHash,
		FirstName:       testUser.FirstName,
		LastName:        testUser.LastName,
		IsActive:        pgtype.Bool{Bool: testUser.IsActive, Valid: true},
		IsEmailVerified: pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: testUser.CreatedAt, Valid: true},
		UpdatedAt:       pgtype.Timestamptz{Time: testUser.UpdatedAt, Valid: true},
	}

	mockQuerier.EXPECT().GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(sqlcUser, nil)

	foundUser, err := repo.GetByID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.ID, foundUser.ID)
	assert.Equal(t, testUser.Email, foundUser.Email)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()

	mockQuerier.EXPECT().GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(sqldb.Users{}, sql.ErrNoRows)

	foundUser, err := repo.GetByID(ctx, userID)
	assert.NoError(t, err) // sql.ErrNoRows should not return an error, but nil user
	assert.Nil(t, foundUser)
}

func TestUserRepository_GetByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()

	mockQuerier.EXPECT().GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(sqldb.Users{}, fmt.Errorf("db error"))

	foundUser, err := repo.GetByID(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Contains(t, err.Error(), "failed to get user by ID")
}

func TestUserRepository_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()
	testUser := &entities.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "John",
		LastName:     "Doe",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sqlcUser := sqldb.Users{
		ID:              pgtype.UUID{Bytes: userID, Valid: true},
		Email:           testUser.Email,
		PasswordHash:    testUser.PasswordHash,
		FirstName:       testUser.FirstName,
		LastName:        testUser.LastName,
		IsActive:        pgtype.Bool{Bool: testUser.IsActive, Valid: true},
		IsEmailVerified: pgtype.Bool{Bool: false, Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: testUser.CreatedAt, Valid: true},
		UpdatedAt:       pgtype.Timestamptz{Time: testUser.UpdatedAt, Valid: true},
	}

	mockQuerier.EXPECT().GetUserByEmail(ctx, testUser.Email).Return(sqlcUser, nil)

	foundUser, err := repo.GetByEmail(ctx, testUser.Email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.ID, foundUser.ID)
	assert.Equal(t, testUser.Email, foundUser.Email)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	testEmail := "nonexistent@example.com"

	mockQuerier.EXPECT().GetUserByEmail(ctx, testEmail).Return(sqldb.Users{}, sql.ErrNoRows)

	foundUser, err := repo.GetByEmail(ctx, testEmail)
	assert.NoError(t, err)
	assert.Nil(t, foundUser)
}

func TestUserRepository_GetByEmail_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	testEmail := "test@example.com"

	mockQuerier.EXPECT().GetUserByEmail(ctx, testEmail).Return(sqldb.Users{}, fmt.Errorf("db error"))

	foundUser, err := repo.GetByEmail(ctx, testEmail)
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Contains(t, err.Error(), "failed to get user by email")
}

func TestUserRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()
	testUser := &entities.User{
		ID:              userID,
		Email:           "updated@example.com",
		PasswordHash:    "newhashedpassword",
		FirstName:       "Jane",
		LastName:        "Doe",
		IsActive:        true,
		IsEmailVerified: true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	sqlcUser := sqldb.Users{
		ID:              pgtype.UUID{Bytes: userID, Valid: true},
		Email:           testUser.Email,
		PasswordHash:    testUser.PasswordHash,
		FirstName:       testUser.FirstName,
		LastName:        testUser.LastName,
		IsActive:        pgtype.Bool{Bool: testUser.IsActive, Valid: true},
		IsEmailVerified: pgtype.Bool{Bool: testUser.IsEmailVerified, Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: testUser.CreatedAt, Valid: true},
		UpdatedAt:       pgtype.Timestamptz{Time: testUser.UpdatedAt, Valid: true},
	}

	updateParams := sqldb.UpdateUserParams{
		ID:              pgtype.UUID{Bytes: testUser.ID, Valid: true},
		Email:           testUser.Email,
		FirstName:       testUser.FirstName,
		LastName:        testUser.LastName,
		IsEmailVerified: pgtype.Bool{Bool: testUser.IsEmailVerified, Valid: true},
	}

	mockQuerier.EXPECT().UpdateUser(ctx, updateParams).Return(sqlcUser, nil)

	updatedUser, err := repo.Update(ctx, testUser)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, testUser.Email, updatedUser.Email)
	assert.Equal(t, testUser.FirstName, updatedUser.FirstName)
}

func TestUserRepository_Update_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()
	testUser := &entities.User{
		ID:              userID,
		Email:           "updated@example.com",
		PasswordHash:    "newhashedpassword",
		FirstName:       "Jane",
		LastName:        "Doe",
		IsActive:        true,
		IsEmailVerified: true,
	}

	updateParams := sqldb.UpdateUserParams{
		ID:              pgtype.UUID{Bytes: testUser.ID, Valid: true},
		Email:           testUser.Email,
		FirstName:       testUser.FirstName,
		LastName:        testUser.LastName,
		IsEmailVerified: pgtype.Bool{Bool: testUser.IsEmailVerified, Valid: true},
	}

	mockQuerier.EXPECT().UpdateUser(ctx, updateParams).Return(sqldb.Users{}, fmt.Errorf("db error"))

	updatedUser, err := repo.Update(ctx, testUser)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Contains(t, err.Error(), "failed to update user")
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()

	mockQuerier.EXPECT().UpdateUserLastLogin(ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(nil)

	err := repo.UpdateLastLogin(ctx, userID)
	assert.NoError(t, err)
}

func TestUserRepository_UpdateLastLogin_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()

	mockQuerier.EXPECT().UpdateUserLastLogin(ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(fmt.Errorf("db error"))

	err := repo.UpdateLastLogin(ctx, userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user last login")
}

func TestUserRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()

	mockQuerier.EXPECT().DeactivateUser(ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(nil)

	err := repo.Delete(ctx, userID)
	assert.NoError(t, err)
}

func TestUserRepository_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	userID := uuid.New()

	mockQuerier.EXPECT().DeactivateUser(ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(fmt.Errorf("db error"))

	err := repo.Delete(ctx, userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete user")
}

func TestUserRepository_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	limit := 10
	offset := 0
	userID1 := uuid.New()
	userID2 := uuid.New()

	sqlcUsers := []sqldb.Users{
		{
			ID:              pgtype.UUID{Bytes: userID1, Valid: true},
			Email:           "test1@example.com",
			PasswordHash:    "hashedpassword1",
			FirstName:       "John",
			LastName:        "Doe",
			IsActive:        pgtype.Bool{Bool: true, Valid: true},
			IsEmailVerified: pgtype.Bool{Bool: false, Valid: true},
			CreatedAt:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:              pgtype.UUID{Bytes: userID2, Valid: true},
			Email:           "test2@example.com",
			PasswordHash:    "hashedpassword2",
			FirstName:       "Jane",
			LastName:        "Smith",
			IsActive:        pgtype.Bool{Bool: true, Valid: true},
			IsEmailVerified: pgtype.Bool{Bool: false, Valid: true},
			CreatedAt:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	listParams := sqldb.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	mockQuerier.EXPECT().ListUsers(ctx, listParams).Return(sqlcUsers, nil)

	users, err := repo.List(ctx, limit, offset)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, userID1, users[0].ID)
	assert.Equal(t, userID2, users[1].ID)
}

func TestUserRepository_List_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_queries.NewMockQuerier(ctrl)
	repo := &userRepository{q: mockQuerier}

	ctx := context.Background()
	limit := 10
	offset := 0

	listParams := sqldb.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	mockQuerier.EXPECT().ListUsers(ctx, listParams).Return(nil, fmt.Errorf("db error"))

	users, err := repo.List(ctx, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "failed to list users")
}