package repository

import (
	"gorm.io/gorm"
	"pojok-baca-api/model"
	"time"
)

type DepositTransactionRepository interface {
	Create(deposit *model.DepositTransaction) (model.DepositTransaction, error)
	UpdateStatus(orderID string, status string, paidAt *time.Time) error
	GetByOrderID(orderID string) (model.DepositTransaction, error)
}

type depositTransactionRepository struct {
	db *gorm.DB
}

func NewDepositTransactionRepository(db *gorm.DB) DepositTransactionRepository {
	return &depositTransactionRepository{db}
}

func (r *depositTransactionRepository) Create(deposit *model.DepositTransaction) (model.DepositTransaction, error) {
	err := r.db.Create(deposit).Error
	return *deposit, err
}

func (r *depositTransactionRepository) UpdateStatus(orderID string, status string, paidAt *time.Time) error {
	return r.db.Model(&model.DepositTransaction{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{"status": status, "paid_at": paidAt}).Error
}

func (r *depositTransactionRepository) GetByOrderID(orderID string) (model.DepositTransaction, error) {
	var tx model.DepositTransaction
	err := r.db.Where("order_id = ?", orderID).First(&tx).Error
	return tx, err
}
