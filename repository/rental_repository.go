package repository

import (
	"gorm.io/gorm"
	"pojok-baca-api/model"
)

type RentalRepository interface {
	Create(rental model.Rental) (model.Rental, error)
	GetByUserID(userID uint) ([]model.Rental, error)
}

type rentalRepository struct {
	db *gorm.DB
}

func NewRentalRepository(db *gorm.DB) RentalRepository {
	return &rentalRepository{db}
}

func (r *rentalRepository) Create(rental model.Rental) (model.Rental, error) {
	err := r.db.Create(&rental).Error
	return rental, err
}

func (r *rentalRepository) GetByUserID(userID uint) ([]model.Rental, error) {
	var rentals []model.Rental
	err := r.db.Preload("Book").Where("user_id = ?", userID).Find(&rentals).Error
	return rentals, err
}
