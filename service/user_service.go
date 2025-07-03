package service

import (
	"pojok-baca-api/model"
	"pojok-baca-api/repository"
)

type UserService interface {
	CreateUser(user model.User) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	UpdateDepositUser(depo int, id uint) (model.User, error)
	GetUserById(id uint) (model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (r *userService) CreateUser(user model.User) (model.User, error) {
	return r.repo.Create(user)
}

func (r *userService) GetUserByEmail(email string) (model.User, error) {
	return r.repo.GetByEmail(email)
}

func (r *userService) UpdateDepositUser(depo int, id uint) (model.User, error) {
	return r.repo.UpdateDeposit(depo, id)
}

func (r *userService) GetUserById(id uint) (model.User, error) {
	return r.repo.GetByID(id)
}
