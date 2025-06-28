package tests

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/user/mocks"
	userUsecase "gofiber-skeleton/internal/user/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userUc := userUsecase.NewUserUsecase(mockRepo)

	// Test case 1: User found
	expectedUser := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
	mockRepo.EXPECT().GetUser(gomock.Any(), uint(1)).Return(expectedUser, nil).Times(1)

	user, err := userUc.GetUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Test case 2: User not found
	mockRepo.EXPECT().GetUser(gomock.Any(), uint(2)).Return(nil, errors.New("user not found")).Times(1)

	user, err = userUc.GetUser(context.Background(), 2)
	assert.Error(t, err)
	assert.Nil(t, user)
}



