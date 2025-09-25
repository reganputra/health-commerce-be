package models

import "time"

type Feedback struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`
	Comment   string
	Rating    int
	CreatedAt time.Time
}