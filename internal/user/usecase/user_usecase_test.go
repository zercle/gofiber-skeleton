package userusecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/domain/mock"
)

func TestUserUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	usecase := NewUserUseCase(mockUserRepo, "test_secret")

	username := "testuser"
	password := "password123"
	role := "customer"

	t.Run("successful registration", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(nil, errors.New("not found")) // User does not exist
		mockUserRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *domain.User) (*domain.User, error) {
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, username, user.Username)
			assert.NotNil(t, user.PasswordHash)
			assert.Equal(t, role, user.Role)
			assert.False(t, user.CreatedAt.IsZero())
			assert.False(t, user.UpdatedAt.IsZero())
			return user, nil
		})

		user, err := usecase.Register(username, password, role)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)
	})

	t.Run("registration with empty role defaults to customer", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(nil, errors.New("not found")) // User does not exist
		mockUserRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *domain.User) (*domain.User, error) {
			assert.Equal(t, domain.RoleCustomer, user.Role)
			return user, nil
		})

		user, err := usecase.Register(username, password, "")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, domain.RoleCustomer, user.Role)
	})

	t.Run("registration with invalid role defaults to customer", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(nil, errors.New("not found")) // User does not exist
		mockUserRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *domain.User) (*domain.User, error) {
			assert.Equal(t, domain.RoleCustomer, user.Role)
			return user, nil
		})

		user, err := usecase.Register(username, password, "invalid_role")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, domain.RoleCustomer, user.Role)
	})

	t.Run("username already exists", func(t *testing.T) {
		existingUser := &domain.User{Username: username}
		mockUserRepo.EXPECT().GetByUsername(username).Return(existingUser, nil)

		user, err := usecase.Register(username, password, role)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "username already exists")
	})

	t.Run("repository create error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(nil, errors.New("not found"))
		mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil, errors.New("db error"))

		user, err := usecase.Register(username, password, role)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	usecase := NewUserUseCase(mockUserRepo, "test_secret")

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleCustomer,
	}

	t.Run("successful login", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(existingUser, nil)

		token, user, err := usecase.Login(username, password)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, user)
		assert.Equal(t, existingUser.ID, user.ID)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(nil, errors.New("not found"))

		token, user, err := usecase.Login(username, password)
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.EqualError(t, err, "invalid credentials")
	})

	t.Run("incorrect password", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(existingUser, nil)

		token, user, err := usecase.Login(username, "wrong_password")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.EqualError(t, err, "invalid credentials")
	})

	t.Run("repository error on GetByUsername", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(username).Return(nil, errors.New("db error"))

		token, user, err := usecase.Login(username, password)
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
		assert.EqualError(t, err, "invalid credentials") // UseCase should return generic "invalid credentials"
	})
}

func TestUserUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	usecase := NewUserUseCase(mockUserRepo, "test_secret")

	userID := uuid.New().String()
	expectedUser := &domain.User{ID: userID, Username: "testuser"}

	t.Run("successful retrieval by ID", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(userID).Return(expectedUser, nil)

		user, err := usecase.GetByID(userID)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(userID).Return(nil, errors.New("db error"))

		user, err := usecase.GetByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserUseCase_UpdateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	usecase := NewUserUseCase(mockUserRepo, "test_secret")

	userID := uuid.New().String()
	originalUser := &domain.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: "hashedpass",
		Role:         domain.RoleCustomer,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("successful role update", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(userID).Return(originalUser, nil)
		mockUserRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(user *domain.User) error {
			assert.Equal(t, domain.RoleAdmin, user.Role)
			assert.False(t, user.UpdatedAt.IsZero())
			return nil
		})

		err := usecase.UpdateRole(userID, domain.RoleAdmin)
		require.NoError(t, err)
	})

	t.Run("invalid role", func(t *testing.T) {
		err := usecase.UpdateRole(userID, "invalid_role")
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid role")
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(userID).Return(nil, errors.New("not found"))

		err := usecase.UpdateRole(userID, domain.RoleAdmin)
		assert.Error(t, err)
		assert.EqualError(t, err, "not found")
	})

	t.Run("repository update error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(userID).Return(originalUser, nil)
		mockUserRepo.EXPECT().Update(gomock.Any()).Return(errors.New("db error"))

		err := usecase.UpdateRole(userID, domain.RoleAdmin)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})
}
