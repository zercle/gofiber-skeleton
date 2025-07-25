// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecases/url_repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/usecases/url_repository.go -destination=mocks/url_repository_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	entities "gofiber-skeleton/internal/entities"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockURLRepository is a mock of URLRepository interface.
type MockURLRepository struct {
	ctrl     *gomock.Controller
	recorder *MockURLRepositoryMockRecorder
	isgomock struct{}
}

// MockURLRepositoryMockRecorder is the mock recorder for MockURLRepository.
type MockURLRepositoryMockRecorder struct {
	mock *MockURLRepository
}

// NewMockURLRepository creates a new mock instance.
func NewMockURLRepository(ctrl *gomock.Controller) *MockURLRepository {
	mock := &MockURLRepository{ctrl: ctrl}
	mock.recorder = &MockURLRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLRepository) EXPECT() *MockURLRepositoryMockRecorder {
	return m.recorder
}

// CreateURL mocks base method.
func (m *MockURLRepository) CreateURL(ctx context.Context, url *entities.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURL", ctx, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateURL indicates an expected call of CreateURL.
func (mr *MockURLRepositoryMockRecorder) CreateURL(ctx, url any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURL", reflect.TypeOf((*MockURLRepository)(nil).CreateURL), ctx, url)
}

// DeleteURL mocks base method.
func (m *MockURLRepository) DeleteURL(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteURL", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteURL indicates an expected call of DeleteURL.
func (mr *MockURLRepositoryMockRecorder) DeleteURL(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteURL", reflect.TypeOf((*MockURLRepository)(nil).DeleteURL), ctx, id)
}

// GetURLByShortCode mocks base method.
func (m *MockURLRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByShortCode", ctx, shortCode)
	ret0, _ := ret[0].(*entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLByShortCode indicates an expected call of GetURLByShortCode.
func (mr *MockURLRepositoryMockRecorder) GetURLByShortCode(ctx, shortCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByShortCode", reflect.TypeOf((*MockURLRepository)(nil).GetURLByShortCode), ctx, shortCode)
}

// GetURLsByUserID mocks base method.
func (m *MockURLRepository) GetURLsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLsByUserID", ctx, userID)
	ret0, _ := ret[0].([]*entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLsByUserID indicates an expected call of GetURLsByUserID.
func (mr *MockURLRepositoryMockRecorder) GetURLsByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLsByUserID", reflect.TypeOf((*MockURLRepository)(nil).GetURLsByUserID), ctx, userID)
}

// UpdateURL mocks base method.
func (m *MockURLRepository) UpdateURL(ctx context.Context, url *entities.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateURL", ctx, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateURL indicates an expected call of UpdateURL.
func (mr *MockURLRepositoryMockRecorder) UpdateURL(ctx, url any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateURL", reflect.TypeOf((*MockURLRepository)(nil).UpdateURL), ctx, url)
}
