package repository

import (
	"pojok-baca-api/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct{
	mock.Mock
}

func (m *UserRepositoryMock) Create(user model.User) (model.User, error){
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}