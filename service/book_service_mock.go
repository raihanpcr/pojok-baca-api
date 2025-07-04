package service

import (
	"pojok-baca-api/dto"
	"pojok-baca-api/model"

	"github.com/stretchr/testify/mock"
)

type BookServiceMock struct {
	mock.Mock
}

func (m *BookServiceMock) GetBooks() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *BookServiceMock) Create(book model.Book) (model.Book, error) {
	args := m.Called(book)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *BookServiceMock) GetBookByID(id uint) (model.Book, error) {
	args := m.Called(id)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *BookServiceMock) DeleteBookByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *BookServiceMock) UpdateBookByID(req dto.UpdateBookRequest, id uint) (model.Book, error) {
	args := m.Called(req, id)
	return args.Get(0).(model.Book), args.Error(1)
}