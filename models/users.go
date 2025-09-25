package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique"`
	Password     string
	Email        string `gorm:"unique"`
	Dob          time.Time
	Gender       string // ENUM('M', 'F')
	Address      string
	City         string
	ContactNumber string
	PaypalID     string
	Role         string // ENUM('admin', 'customer')
	CreatedAt    time.Time
	UpdatedAt    time.Time
}