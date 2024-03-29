package mocks

import (
	mock "github.com/stretchr/testify/mock"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type BookRepository struct {
	mock.Mock
}

func (_m *BookRepository) GetBook(bookID uint) (book models.Book, err error) {
	ret := _m.Called(bookID)

	var r0 models.Book
	if rf, ok := ret.Get(0).(func(uint) models.Book); ok {
		r0 = rf(bookID)
	} else {
		r0 = ret.Get(0).(models.Book)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(bookID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *BookRepository) GetBooks(criteria models.Book) (books []models.Book, err error) {
	ret := _m.Called(criteria)

	var r0 []models.Book
	if rf, ok := ret.Get(0).(func(models.Book) []models.Book); ok {
		r0 = rf(criteria)
	} else {
		r0 = ret.Get(0).([]models.Book)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Book) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *BookRepository) CreateBook(book *models.Book) (err error) {
	ret := _m.Called(book)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Book) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

func (_m *BookRepository) EditBook(bookID uint, book models.Book) (err error) {
	ret := _m.Called(bookID, book)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, models.Book) error); ok {
		r0 = rf(bookID, book)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

func (_m *BookRepository) DeleteBook(bookID uint) (err error) {
	ret := _m.Called(bookID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(bookID)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}
