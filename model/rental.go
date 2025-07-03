package model

import (
	"gorm.io/gorm"
	"time"
)

type Rental struct {
	gorm.Model
	UserID     uint      `gorm:"not null"`
	BookID     uint      `gorm:"not null"`
	RentDate   time.Time `gorm:"not null"`
	ReturnDate *time.Time
	Status     string `gorm:"not null"`
	User       User
	Book       Book
}
