package service

import (
	"pojok-baca-api/model"

	"github.com/stretchr/testify/mock"
)

type RentalServiceMock struct {
	mock.Mock
}

func (m *RentalServiceMock) CreateRental(rental model.Rental) (model.Rental, error) {
	args := m.Called(rental)
	return args.Get(0).(model.Rental), args.Error(1)
}

func (m *RentalServiceMock) GetRentalByUserID(userID uint) ([]model.Rental, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Rental), args.Error(1)
}