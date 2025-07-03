package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name                string `gorm:"not null"`
	Email               string `gorm:"unique; not null"`
	Password            string `gorm:"not null"`
	Deposit             *int
	Role                string               `gorm:"not null"`
	Rental              []Rental             `gorm:"foreignKey:UserID"`
	DepositTransactions []DepositTransaction `gorm:"foreignKey:UserID"`
}
