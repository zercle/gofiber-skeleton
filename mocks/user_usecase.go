package mocks

import (
	mock "github.com/stretchr/testify/mock"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type UserUsecase struct {
	mock.Mock
}

func (_m *UserUsecase) GetUser(userId string) (user models.User, err error) {
	ret := _m.Called(userId)

	var r0 models.User
	if rf, ok := ret.Get(0).(func(string) models.User); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

func (_m *UserUsecase) GetUsers(criteria models.User) (users []models.User, err error) {
	ret := _m.Called(criteria)

	var r0 []models.User
	if rf, ok := ret.Get(0).(func(models.User) []models.User); ok {
		r0 = rf(criteria)
	} else {
		r0 = ret.Get(0).([]models.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

func (_m *UserUsecase) CreateUser(user *models.User) (err error) {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

func (_m *UserUsecase) EditUser(userId string, user models.User) (err error) {
	ret := _m.Called(userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

func (_m *UserUsecase) DeleteUser(userId string) (err error) {
	ret := _m.Called(userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}
