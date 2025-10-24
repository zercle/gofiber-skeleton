package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/zercle/template-go-fiber/internal/domains"
	"github.com/zercle/template-go-fiber/internal/usecases"
	"github.com/zercle/template-go-fiber/test/unit/mocks"
)

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	// Expect GetByEmail to be called to check for duplicates
	mockRepo.EXPECT().
		GetByEmail(gomock.Any(), "newuser@example.com").
		Return(nil, nil).
		Times(1)

	// Expect Create to be called with user data
	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	input := &domains.RegisterUserInput{
		Email:    "newuser@example.com",
		Password: "password123",
	}

	user, err := usecase.RegisterUser(context.Background(), input)

	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Email != "newuser@example.com" {
		t.Errorf("expected email %s, got %s", "newuser@example.com", user.Email)
	}

	if !user.IsActive {
		t.Error("expected user to be active")
	}
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	existingUser := &domains.User{
		ID:    "123",
		Email: "existing@example.com",
	}

	// Expect GetByEmail to return existing user
	mockRepo.EXPECT().
		GetByEmail(gomock.Any(), "existing@example.com").
		Return(existingUser, nil).
		Times(1)

	input := &domains.RegisterUserInput{
		Email:    "existing@example.com",
		Password: "password123",
	}

	user, err := usecase.RegisterUser(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestRegisterUser_MissingEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	input := &domains.RegisterUserInput{
		Email:    "",
		Password: "password123",
	}

	user, err := usecase.RegisterUser(context.Background(), input)

	if err == nil {
		t.Fatal("expected error for missing email")
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestRegisterUser_MissingPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	input := &domains.RegisterUserInput{
		Email:    "user@example.com",
		Password: "",
	}

	user, err := usecase.RegisterUser(context.Background(), input)

	if err == nil {
		t.Fatal("expected error for missing password")
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	expectedUser := &domains.User{
		ID:       "123",
		Email:    "user@example.com",
		IsActive: true,
		CreatedAt: time.Now(),
	}

	mockRepo.EXPECT().
		GetByID(gomock.Any(), "123").
		Return(expectedUser, nil).
		Times(1)

	user, err := usecase.GetUserByID(context.Background(), "123")

	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if user.ID != "123" {
		t.Errorf("expected ID %s, got %s", "123", user.ID)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, nil).
		Times(1)

	user, err := usecase.GetUserByID(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error for user not found")
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	expectedUser := &domains.User{
		ID:       "123",
		Email:    "user@example.com",
		IsActive: true,
	}

	mockRepo.EXPECT().
		GetByEmail(gomock.Any(), "user@example.com").
		Return(expectedUser, nil).
		Times(1)

	user, err := usecase.GetUserByEmail(context.Background(), "user@example.com")

	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}

	if user.Email != "user@example.com" {
		t.Errorf("expected email %s, got %s", "user@example.com", user.Email)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	existingUser := &domains.User{
		ID:       "123",
		Email:    "user@example.com",
		IsActive: true,
	}

	mockRepo.EXPECT().
		GetByID(gomock.Any(), "123").
		Return(existingUser, nil).
		Times(1)

	newEmail := "newemail@example.com"
	mockRepo.EXPECT().
		GetByEmail(gomock.Any(), newEmail).
		Return(nil, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	input := &domains.UpdateUserInput{
		Email: &newEmail,
	}

	user, err := usecase.UpdateUser(context.Background(), "123", input)

	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if user.Email != newEmail {
		t.Errorf("expected email %s, got %s", newEmail, user.Email)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, nil).
		Times(1)

	input := &domains.UpdateUserInput{}

	user, err := usecase.UpdateUser(context.Background(), "nonexistent", input)

	if err == nil {
		t.Fatal("expected error for user not found")
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	existingUser := &domains.User{
		ID:    "123",
		Email: "user@example.com",
	}

	mockRepo.EXPECT().
		GetByID(gomock.Any(), "123").
		Return(existingUser, nil).
		Times(1)

	mockRepo.EXPECT().
		Delete(gomock.Any(), "123").
		Return(nil).
		Times(1)

	err := usecase.DeleteUser(context.Background(), "123")

	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), "nonexistent").
		Return(nil, nil).
		Times(1)

	err := usecase.DeleteUser(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error for user not found")
	}
}

func TestListUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	expectedUsers := []*domains.User{
		{ID: "1", Email: "user1@example.com"},
		{ID: "2", Email: "user2@example.com"},
	}

	mockRepo.EXPECT().
		List(gomock.Any(), int32(10), int32(0)).
		Return(expectedUsers, nil).
		Times(1)

	users, err := usecase.ListUsers(context.Background(), 10, 0)

	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestListUsers_PaginationDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := usecases.NewUserUsecase(mockRepo)

	mockRepo.EXPECT().
		List(gomock.Any(), int32(10), int32(0)).
		Return([]*domains.User{}, nil).
		Times(1)

	// Test that default limit is applied
	_, err := usecase.ListUsers(context.Background(), 0, -1)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
}
