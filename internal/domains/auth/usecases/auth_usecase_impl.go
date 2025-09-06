package usecases

import (
	"context"
	"fmt"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/entities"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/models"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/repositories"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
	"github.com/zercle/gofiber-skeleton/pkg/utils"
)

type authUsecaseImpl struct {
	userRepo   repositories.UserRepository
	jwtManager *utils.JWTManager
}

func NewAuthUsecase(userRepo repositories.UserRepository, jwtManager *utils.JWTManager) AuthUsecase {
	return &authUsecaseImpl{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (u *authUsecaseImpl) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	
	if exists {
		return nil, types.ErrEmailExists
	}

	user, err := entities.NewUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := u.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  models.NewUserResponse(user),
	}, nil
}

func (u *authUsecaseImpl) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == types.ErrUserNotFound {
			return nil, types.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, types.ErrUnauthorized
	}

	if !user.CheckPassword(req.Password) {
		return nil, types.ErrInvalidCredentials
	}

	token, err := u.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  models.NewUserResponse(user),
	}, nil
}

func (u *authUsecaseImpl) GetProfile(ctx context.Context, userID string) (*models.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := models.NewUserResponse(user)
	return &response, nil
}

func (u *authUsecaseImpl) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := models.NewUserResponse(user)
	return &response, nil
}

func (u *authUsecaseImpl) ChangePassword(ctx context.Context, userID string, req models.ChangePasswordRequest) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.CheckPassword(req.CurrentPassword) {
		return types.ErrInvalidCredentials
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *authUsecaseImpl) ListUsers(ctx context.Context, page, pageSize int) (*models.UsersListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, err := u.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	response := models.NewUsersListResponse(users, total, page, pageSize)
	return &response, nil
}

func (u *authUsecaseImpl) GetUser(ctx context.Context, userID string) (*models.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := models.NewUserResponse(user)
	return &response, nil
}

func (u *authUsecaseImpl) DeactivateUser(ctx context.Context, userID string) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Deactivate()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *authUsecaseImpl) ActivateUser(ctx context.Context, userID string) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Activate()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (u *authUsecaseImpl) DeleteUser(ctx context.Context, userID string) error {
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}