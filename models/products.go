package models

import (
	"time"
)

type Product struct {
	ID          uint `gorm:"primaryKey"`
	CategoryID  uint
	Category    Category `gorm:"foreignKey:CategoryID"`
	Name        string
	Description string
	Price       float64
	Stock       int
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}