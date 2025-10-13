package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrEmailAlreadyUsed = errors.New("email already used")
)

type UserUsecase interface {
	Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.User, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, req *entity.UpdateUserRequest) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
	jwtSecret string
}

func NewUserUsecase(userRepo repository.UserRepository, jwtSecret string) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.User, error) {
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}

	passwordHash, err := u.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := entity.NewUser(req.Email, passwordHash, req.FullName)

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !u.verifyPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	token, err := u.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entity.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *userUsecase) Update(ctx context.Context, id uuid.UUID, req *entity.UpdateUserRequest) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Email != nil {
		existingUser, err := u.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailAlreadyUsed
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		passwordHash, err := u.hashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = passwordHash
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (u *userUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (u *userUsecase) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

func (u *userUsecase) hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := randRead(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return fmt.Sprintf("%x:%x", salt, hash), nil
}

func (u *userUsecase) verifyPassword(password, hash string) bool {
	var salt, key []byte
	if _, err := fmt.Sscanf(hash, "%x:%x", &salt, &key); err != nil {
		return false
	}

	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return string(key) == string(computedHash)
}

func (u *userUsecase) generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}

func randRead(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	for i := range b {
		b[i] = byte(time.Now().UnixNano() & 0xff)
	}
	return len(b), nil
}