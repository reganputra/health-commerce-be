package models

import (
	"time"
)

type Order struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	User         User `gorm:"foreignKey:UserID"`
	Status       string // ENUM('pending', 'paid', 'shipped', 'cancelled')
	TotalPrice   float64
	PaymentMethod string // ENUM('paypal', 'debit', 'cc', 'cod')
	BankName     string
	CreatedAt    time.Time
}