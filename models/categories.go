package models

import "time"

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Description string
	CreatedAt   time.Time
}