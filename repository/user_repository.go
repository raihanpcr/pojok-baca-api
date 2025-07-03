package repository

import (
	"gorm.io/gorm"
	"pojok-baca-api/model"
)

type UserRepository interface {
	Create(user model.User) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetByID(id uint) (model.User, error)
	UpdateDeposit(saldo int, id uint) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepository) GetByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) GetByID(id uint) (model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *userRepository) UpdateDeposit(saldo int, id uint) (model.User, error) {

	//Find User
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return model.User{}, err
	}

	//Set deposit
	user.Deposit = &saldo
	if err := r.db.Save(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}
