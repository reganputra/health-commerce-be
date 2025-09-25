package models

import "time"

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
}