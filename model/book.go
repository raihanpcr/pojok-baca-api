package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name       string   `gorm:"not null"`
	Stok       int      `gorm:"not null"`
	RentalCost int      `gorm:"not null"`
	Category   string   `gorm:"not null"`
	Rental     []Rental `gorm:"foreignKey:BookID"`
}
