package model

import (
	"gorm.io/gorm"
	"time"
)

type DepositTransaction struct {
	gorm.Model
	UserID     uint
	OrderID    string
	PaymentRef string 
	Deposit    int
	Status     string
	PaidAt     *time.Time
}
