package mocks

import "gofiber-skeleton/internal/domain"

type BookingRepoMock struct {
	CreateFn      func(b *domain.Booking) error
	GetByIDFn     func(id uint) (*domain.Booking, error)
	UpdateFn      func(b *domain.Booking) error
	DeleteFn      func(id uint) error
	ListByUserFn  func(userID uint) ([]domain.Booking, error)
}

func (m *BookingRepoMock) Create(b *domain.Booking) error {
	if m.CreateFn != nil {
		return m.CreateFn(b)
	}
	return nil
}

func (m *BookingRepoMock) GetByID(id uint) (*domain.Booking, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return &domain.Booking{}, nil
}

func (m *BookingRepoMock) Update(b *domain.Booking) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(b)
	}
	return nil
}

func (m *BookingRepoMock) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func (m *BookingRepoMock) ListByUser(userID uint) ([]domain.Booking, error) {
	if m.ListByUserFn != nil {
		return m.ListByUserFn(userID)
	}
	return nil, nil
}