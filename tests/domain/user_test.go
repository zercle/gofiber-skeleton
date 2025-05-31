package domain_test

import (
	"errors"
	"testing"
	"time"

	"gofiber-skeleton/internal/domain"
)

// MockUserRepo implements domain.UserRepository for testing
type MockUserRepo struct {
	CreateFn       func(user *domain.User) error
	GetByIDFn      func(id uint) (*domain.User, error)
	GetByUsernameFn func(username string) (*domain.User, error)
	GetByEmailFn   func(email string) (*domain.User, error)
	UpdateFn       func(user *domain.User) error
	DeleteFn       func(id uint) error
}

func (m *MockUserRepo) Create(user *domain.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(user)
	}
	return nil
}
func (m *MockUserRepo) GetByID(id uint) (*domain.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, errors.New("not implemented")
}
func (m *MockUserRepo) GetByUsername(username string) (*domain.User, error) {
	if m.GetByUsernameFn != nil {
		return m.GetByUsernameFn(username)
	}
	return nil, errors.New("not implemented")
}
func (m *MockUserRepo) GetByEmail(email string) (*domain.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(email)
	}
	return nil, errors.New("not implemented")
}
func (m *MockUserRepo) Update(user *domain.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(user)
	}
	return nil
}
func (m *MockUserRepo) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func TestUserStructFields(t *testing.T) {
	now := time.Now()
	user := domain.User{
		ID:        1,
		Username:  "alice",
		Email:     "alice@example.com",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	if user.Username != "alice" {
		t.Errorf("expected Username 'alice', got %s", user.Username)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected Email 'alice@example.com', got %s", user.Email)
	}
	if user.Password != "hashedpassword" {
		t.Errorf("expected Password 'hashedpassword', got %s", user.Password)
	}
	if !user.CreatedAt.Equal(now) || !user.UpdatedAt.Equal(now) {
		t.Errorf("expected CreatedAt/UpdatedAt to match now")
	}
}

func TestUserRepositoryMock(t *testing.T) {
	mock := &MockUserRepo{
		CreateFn: func(u *domain.User) error {
			if u.Username == "" {
				return errors.New("missing username")
			}
			return nil
		},
		GetByIDFn: func(id uint) (*domain.User, error) {
			if id == 42 {
				return &domain.User{ID: 42, Username: "bob"}, nil
			}
			return nil, errors.New("not found")
		},
		GetByUsernameFn: func(username string) (*domain.User, error) {
			if username == "alice" {
				return &domain.User{ID: 1, Username: "alice"}, nil
			}
			return nil, errors.New("not found")
		},
		GetByEmailFn: func(email string) (*domain.User, error) {
			if email == "bob@example.com" {
				return &domain.User{ID: 2, Email: "bob@example.com"}, nil
			}
			return nil, errors.New("not found")
		},
		UpdateFn: func(u *domain.User) error {
			if u.ID == 0 {
				return errors.New("missing id")
			}
			return nil
		},
		DeleteFn: func(id uint) error {
			if id == 0 {
				return errors.New("missing id")
			}
			return nil
		},
	}

	// Create
	err := mock.Create(&domain.User{Username: ""})
	if err == nil {
		t.Error("expected error for missing username")
	}
	err = mock.Create(&domain.User{Username: "alice"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// GetByID
	u, err := mock.GetByID(42)
	if err != nil || u.Username != "bob" {
		t.Errorf("expected bob, got %v, err %v", u, err)
	}
	_, err = mock.GetByID(99)
	if err == nil {
		t.Error("expected error for unknown id")
	}

	// GetByUsername
	u, err = mock.GetByUsername("alice")
	if err != nil || u.ID != 1 {
		t.Errorf("expected alice, got %v, err %v", u, err)
	}
	_, err = mock.GetByUsername("nobody")
	if err == nil {
		t.Error("expected error for unknown username")
	}

	// GetByEmail
	u, err = mock.GetByEmail("bob@example.com")
	if err != nil || u.ID != 2 {
		t.Errorf("expected bob@example.com, got %v, err %v", u, err)
	}
	_, err = mock.GetByEmail("none@example.com")
	if err == nil {
		t.Error("expected error for unknown email")
	}

	// Update
	err = mock.Update(&domain.User{ID: 0})
	if err == nil {
		t.Error("expected error for missing id")
	}
	err = mock.Update(&domain.User{ID: 5})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Delete
	err = mock.Delete(0)
	if err == nil {
		t.Error("expected error for missing id")
	}
	err = mock.Delete(7)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
