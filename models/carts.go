package models

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"column:user_id;not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
