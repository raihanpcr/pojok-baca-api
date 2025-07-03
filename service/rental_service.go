package service

import (
	"errors"
	"pojok-baca-api/model"
	"pojok-baca-api/repository"
)

type RentalService interface {
	CreateRental(rental model.Rental) (model.Rental, error)
	GetRentalByUserID(userID uint) ([]model.Rental, error)
}

type rentalService struct {
	repo repository.RentalRepository
}

func NewRentalService(repo repository.RentalRepository) RentalService {
	return &rentalService{repo: repo}
}

func (s *rentalService) CreateRental(rent model.Rental) (model.Rental, error) {
	if rent.BookID == 0 {
		return model.Rental{}, errors.New("BookID is required")
	}

	return s.repo.Create(rent)
}

func (s *rentalService) GetRentalByUserID(userID uint) ([]model.Rental, error) {
	return s.repo.GetByUserID(userID)
}
