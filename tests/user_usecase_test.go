package tests

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/user/mocks"
	userUsecase "gofiber-skeleton/internal/user/usecase"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtSecret := "secretForTest"
	authSrv := auth.NewJWTService(jwtSecret, 24*time.Hour)
	mockRepo := mocks.NewMockUserRepository(ctrl)
	userUc := userUsecase.NewUserUsecase(mockRepo, authSrv)

	existingUserID, _ := uuid.NewV7()
	wrongUserID, _ := uuid.NewV7()

	// Test case 1: User found
	expectedUser := &domain.User{ID: existingUserID, Name: "testuser", Email: "test@example.com"}
	mockRepo.EXPECT().GetUser(gomock.Any(), uint(1)).Return(expectedUser, nil).Times(1)

	user, err := userUc.GetUser(context.Background(), existingUserID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Test case 2: User not found
	mockRepo.EXPECT().GetUser(gomock.Any(), wrongUserID).Return(nil, errors.New("user not found")).Times(1)

	user, err = userUc.GetUser(context.Background(), wrongUserID)
	assert.Error(t, err)
	assert.Nil(t, user)
}
