package service

import (
	"pojok-baca-api/model"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) CreateUser(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserServiceMock) GetUserByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserServiceMock) UpdateDepositUser(depo int, id uint) (model.User, error) {
	args := m.Called(depo, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserServiceMock) GetUserById(id uint) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}