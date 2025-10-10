package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
)

// mockUserRepository is a simple mock for testing
type mockUserRepository struct {
	users map[string]*entity.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	m.users[user.Email] = user
	return user, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, assert.AnError
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, assert.AnError
	}
	return user, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	m.users[user.Email] = user
	return user, nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return assert.AnError
}

func (m *mockUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	users := make([]*entity.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *mockUserRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.users)), nil
}

func (m *mockUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	for _, user := range m.users {
		if user.ID == id {
			user.PasswordHash = passwordHash
			return nil
		}
	}
	return assert.AnError
}

func (m *mockUserRepository) Verify(ctx context.Context, id uuid.UUID) error {
	for _, user := range m.users {
		if user.ID == id {
			user.IsVerified = true
			return nil
		}
	}
	return assert.AnError
}

func getTestConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test-secret-key",
			AccessExpiration: 15 * time.Minute,
		},
	}
}

func TestRegister(t *testing.T) {
	repo := newMockUserRepository()
	cfg := getTestConfig()
	usecase := NewAuthUsecase(repo, cfg)

	tests := []struct {
		name    string
		req     *entity.RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful registration",
			req: &entity.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &entity.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User 2",
			},
			wantErr: true,
			errMsg:  "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := usecase.Register(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.req.Email, user.Email)
				assert.Equal(t, tt.req.FullName, user.FullName)
				assert.NotEmpty(t, user.ID)
				assert.True(t, user.IsActive)
				assert.False(t, user.IsVerified)

				// Verify password was hashed
				err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.req.Password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	repo := newMockUserRepository()
	cfg := getTestConfig()
	usecase := NewAuthUsecase(repo, cfg)

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		IsActive:     true,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	repo.users[testUser.Email] = testUser

	tests := []struct {
		name    string
		req     *entity.LoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful login",
			req: &entity.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "wrong password",
			req: &entity.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "user not found",
			req: &entity.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := usecase.Login(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.Equal(t, "Bearer", resp.TokenType)
				assert.Greater(t, resp.ExpiresIn, int64(0))
				assert.Equal(t, testUser.Email, resp.User.Email)
			}
		})
	}
}

func TestLoginInactiveUser(t *testing.T) {
	repo := newMockUserRepository()
	cfg := getTestConfig()
	usecase := NewAuthUsecase(repo, cfg)

	// Create an inactive user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &entity.User{
		ID:           uuid.New(),
		Email:        "inactive@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Inactive User",
		IsActive:     false,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	repo.users[testUser.Email] = testUser

	resp, err := usecase.Login(context.Background(), &entity.LoginRequest{
		Email:    "inactive@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deactivated")
	assert.Nil(t, resp)
}

func TestVerifyToken(t *testing.T) {
	repo := newMockUserRepository()
	cfg := getTestConfig()
	usecase := NewAuthUsecase(repo, cfg)

	userID := uuid.New()

	tests := []struct {
		name      string
		setupFunc func() string
		wantErr   bool
		checkID   bool
	}{
		{
			name: "valid token",
			setupFunc: func() string {
				claims := jwt.MapClaims{
					"user_id": userID.String(),
					"exp":     time.Now().Add(15 * time.Minute).Unix(),
					"iat":     time.Now().Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.JWT.Secret))
				return tokenString
			},
			wantErr: false,
			checkID: true,
		},
		{
			name: "expired token",
			setupFunc: func() string {
				claims := jwt.MapClaims{
					"user_id": userID.String(),
					"exp":     time.Now().Add(-1 * time.Hour).Unix(),
					"iat":     time.Now().Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(cfg.JWT.Secret))
				return tokenString
			},
			wantErr: true,
		},
		{
			name: "invalid signature",
			setupFunc: func() string {
				claims := jwt.MapClaims{
					"user_id": userID.String(),
					"exp":     time.Now().Add(15 * time.Minute).Unix(),
					"iat":     time.Now().Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("wrong-secret"))
				return tokenString
			},
			wantErr: true,
		},
		{
			name: "malformed token",
			setupFunc: func() string {
				return "malformed.token.string"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupFunc()
			parsedUserID, err := usecase.VerifyToken(token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkID {
					assert.Equal(t, userID, parsedUserID)
				}
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	repo := newMockUserRepository()
	cfg := getTestConfig()
	usecase := NewAuthUsecase(repo, cfg)

	testUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		FullName:     "Old Name",
		IsActive:     true,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	repo.users[testUser.Email] = testUser

	req := &entity.UpdateProfileRequest{
		FullName: "New Name",
	}

	updatedUser, err := usecase.UpdateProfile(context.Background(), testUser.ID, req)

	require.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "New Name", updatedUser.FullName)
}

func TestChangePassword(t *testing.T) {
	repo := newMockUserRepository()
	cfg := getTestConfig()
	usecase := NewAuthUsecase(repo, cfg)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	testUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		IsActive:     true,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	repo.users[testUser.Email] = testUser

	tests := []struct {
		name    string
		req     *entity.ChangePasswordRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful password change",
			req: &entity.ChangePasswordRequest{
				OldPassword: "oldpassword",
				NewPassword: "newpassword123",
			},
			wantErr: false,
		},
		{
			name: "wrong old password",
			req: &entity.ChangePasswordRequest{
				OldPassword: "wrongpassword",
				NewPassword: "newpassword123",
			},
			wantErr: true,
			errMsg:  "invalid old password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.ChangePassword(context.Background(), testUser.ID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)

				// Verify new password was set
				user, _ := repo.GetByID(context.Background(), testUser.ID)
				err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.req.NewPassword))
				assert.NoError(t, err)
			}
		})
	}
}
