package usecase_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"gofiber-skeleton/internal/domain"
	"gofiber-skeleton/internal/usecase"
	"gofiber-skeleton/tests/mocks"
)

func TestCreateBooking(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-time.Hour)
	later := now.Add(time.Hour)

	tests := []struct {
		name     string
		input    *domain.Booking
		repoErr  error
		wantErr  error
	}{
		{"nil booking", nil, nil, usecase.ErrInvalidBooking},
		{"missing user", &domain.Booking{UserID: 0, ItemID: 1, StartDate: earlier, EndDate: later}, nil, usecase.ErrInvalidBooking},
		{"missing item", &domain.Booking{UserID: 1, ItemID: 0, StartDate: earlier, EndDate: later}, nil, usecase.ErrInvalidBooking},
		{"date conflict", &domain.Booking{UserID: 1, ItemID: 2, StartDate: later, EndDate: earlier}, nil, usecase.ErrBookingDateConflict},
		{"repo error", &domain.Booking{UserID: 1, ItemID: 2, StartDate: earlier, EndDate: later}, errors.New("fail"), errors.New("fail")},
		{"success", &domain.Booking{UserID: 1, ItemID: 2, StartDate: earlier, EndDate: later}, nil, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mocks.BookingRepoMock{
				CreateFn: func(b *domain.Booking) error {
					return tc.repoErr
				},
			}
			svc := usecase.NewBookingUsecase(mockRepo)
			err := svc.CreateBooking(tc.input)
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
			}
		})
	}
}

func TestGetBooking(t *testing.T) {
	expected := &domain.Booking{ID: 5}
	tests := []struct {
		name       string
		inputID    uint
		repoValue  *domain.Booking
		repoErr    error
		wantValue  *domain.Booking
		wantErr    error
	}{
		{"invalid id", 0, nil, nil, nil, usecase.ErrInvalidBooking},
		{"not found", 1, nil, errors.New("db"), nil, usecase.ErrBookingNotFound},
		{"success", 5, expected, nil, expected, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mocks.BookingRepoMock{
				GetByIDFn: func(id uint) (*domain.Booking, error) {
					return tc.repoValue, tc.repoErr
				},
			}
			svc := usecase.NewBookingUsecase(mockRepo)
			got, err := svc.GetBooking(tc.inputID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if !reflect.DeepEqual(got, tc.wantValue) {
					t.Errorf("expected %v, got %v", tc.wantValue, got)
				}
			}
		})
	}
}

func TestUpdateBooking(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-time.Hour)
	later := now.Add(time.Hour)

	tests := []struct {
		name     string
		input    *domain.Booking
		repoErr  error
		wantErr  error
	}{
		{"nil booking", nil, nil, usecase.ErrInvalidBooking},
		{"no id", &domain.Booking{ID: 0, UserID: 1, ItemID: 2, StartDate: earlier, EndDate: later}, nil, usecase.ErrInvalidBooking},
		{"missing user", &domain.Booking{ID: 1, UserID: 0, ItemID: 2, StartDate: earlier, EndDate: later}, nil, usecase.ErrInvalidBooking},
		{"missing item", &domain.Booking{ID: 1, UserID: 1, ItemID: 0, StartDate: earlier, EndDate: later}, nil, usecase.ErrInvalidBooking},
		{"date conflict", &domain.Booking{ID: 1, UserID: 1, ItemID: 2, StartDate: later, EndDate: earlier}, nil, usecase.ErrBookingDateConflict},
		{"repo error", &domain.Booking{ID: 1, UserID: 1, ItemID: 2, StartDate: earlier, EndDate: later}, errors.New("update fail"), errors.New("update fail")},
		{"success", &domain.Booking{ID: 1, UserID: 1, ItemID: 2, StartDate: earlier, EndDate: later}, nil, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mocks.BookingRepoMock{
				UpdateFn: func(b *domain.Booking) error {
					return tc.repoErr
				},
			}
			svc := usecase.NewBookingUsecase(mockRepo)
			err := svc.UpdateBooking(tc.input)
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
			}
		})
	}
}

func TestDeleteBooking(t *testing.T) {
	tests := []struct {
		name    string
		inputID uint
		repoErr error
		wantErr error
	}{
		{"invalid id", 0, nil, usecase.ErrInvalidBooking},
		{"repo error", 1, errors.New("del fail"), errors.New("del fail")},
		{"success", 2, nil, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mocks.BookingRepoMock{
				DeleteFn: func(id uint) error {
					return tc.repoErr
				},
			}
			svc := usecase.NewBookingUsecase(mockRepo)
			err := svc.DeleteBooking(tc.inputID)
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
			}
		})
	}
}

func TestListBookingsByUser(t *testing.T) {
	expected := []domain.Booking{{ID: 10}, {ID: 20}}
	tests := []struct {
		name     string
		inputID  uint
		repoList []domain.Booking
		repoErr  error
		wantList []domain.Booking
		wantErr  error
	}{
		{"invalid user", 0, nil, nil, nil, usecase.ErrInvalidBooking},
		{"repo error", 1, nil, errors.New("list fail"), nil, errors.New("list fail")},
		{"success", 2, expected, nil, expected, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mocks.BookingRepoMock{
				ListByUserFn: func(userID uint) ([]domain.Booking, error) {
					return tc.repoList, tc.repoErr
				},
			}
			svc := usecase.NewBookingUsecase(mockRepo)
			list, err := svc.ListBookingsByUser(tc.inputID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if !reflect.DeepEqual(list, tc.wantList) {
					t.Errorf("expected %v, got %v", tc.wantList, list)
				}
			}
		})
	}
}