package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/internal/testutil"
	"github.com/zercle/gofiber-skeleton/internal/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/user/repository/mocks"
	"github.com/zercle/gofiber-skeleton/internal/user/usecase"
)

// MockUserRepository is a mock for testing
//go:generate mockgen -destination=mocks/repository_mock.go -package=mocks github.com/zercle/gofiber-skeleton/internal/user/repository UserRepository

func TestAuthUsecase_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockRepo, testutil.ValidJWTSecret())

	ctx := context.Background()
	username := "newuser"
	email := "newuser@example.com"
	password := "password123"

	// Expect Create to be called with any user and password hash
	mockRepo.EXPECT().
		Create(ctx, gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *entity.User, passwordHash string) error {
			// Verify password was hashed
			err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
			assert.NoError(t, err, "Password should be properly hashed")

			// Verify user fields
			assert.Equal(t, username, user.Username)
			assert.Equal(t, email, user.Email)
			assert.NotEmpty(t, user.ID)

			return nil
		}).
		Times(1)

	result, err := authUsecase.Register(ctx, username, email, password)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, username, result.Username)
	assert.Equal(t, email, result.Email)
	assert.NotEmpty(t, result.ID)
}

func TestAuthUsecase_Register_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockRepo, testutil.ValidJWTSecret())

	ctx := context.Background()
	expectedErr := errors.New("database error")

	mockRepo.EXPECT().
		Create(ctx, gomock.Any(), gomock.Any()).
		Return(expectedErr).
		Times(1)

	result, err := authUsecase.Register(ctx, "user", "user@example.com", "password")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockRepo, testutil.ValidJWTSecret())

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	user := testutil.UserFixture(t, func(u *entity.User) {
		u.Email = email
	})

	// Generate hash for the test password
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	passwordHash := string(passwordHashBytes)

	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(user, passwordHash, nil).
		Times(1)

	token, resultUser, err := authUsecase.Login(ctx, email, password)

	require.NoError(t, err)
	require.NotNil(t, resultUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.ID, resultUser.ID)
	assert.Equal(t, user.Email, resultUser.Email)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockRepo, testutil.ValidJWTSecret())

	ctx := context.Background()
	email := "notfound@example.com"

	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, "", sql.ErrNoRows).
		Times(1)

	token, user, err := authUsecase.Login(ctx, email, "password")

	require.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "not found")
}

func TestAuthUsecase_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockRepo, testutil.ValidJWTSecret())

	ctx := context.Background()
	email := "test@example.com"

	user := testutil.UserFixture(t, func(u *entity.User) {
		u.Email = email
	})

	passwordHash := testutil.PasswordHash() // Hash of "password123"

	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(user, passwordHash, nil).
		Times(1)

	// Try with wrong password
	token, resultUser, err := authUsecase.Login(ctx, email, "wrongpassword")

	require.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, resultUser)
}

func TestAuthUsecase_Login_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockRepo, testutil.ValidJWTSecret())

	ctx := context.Background()
	email := "test@example.com"
	expectedErr := errors.New("database connection error")

	mockRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, "", expectedErr).
		Times(1)

	token, user, err := authUsecase.Login(ctx, email, "password")

	require.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Equal(t, expectedErr, err)
}
