package tests

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/user/mocks"
	userUsecase "gofiber-skeleton/internal/user/usecase"
	"gofiber-skeleton/internal/infra/types"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userUc := userUsecase.NewUserUsecase(mockRepo)

	// Test case 1: User found
	userID := types.NewUUIDv7()
	expectedUser := &domain.User{ID: userID, Username: "testuser", Email: "test@example.com"}
	mockRepo.EXPECT().GetUser(gomock.Any(), userID).Return(expectedUser, nil).Times(1)

	user, err := userUc.GetUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Test case 2: User not found
	notFoundUserID := types.NewUUIDv7()
	mockRepo.EXPECT().GetUser(gomock.Any(), notFoundUserID).Return(nil, errors.New("user not found")).Times(1)

	user, err = userUc.GetUser(context.Background(), notFoundUserID)
	assert.Error(t, err)
	assert.Nil(t, user)
}



