package models

import "time"

type Feedback struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID uint      `json:"productId"`
	Product   Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Comment   string    `json:"comment"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"createdAt"`
}

// FeedbackResponse represents the response structure for feedback with selected fields
type FeedbackResponse struct {
	ID        uint        `json:"id"`
	UserID    uint        `json:"userId"`
	User      UserInfo    `json:"user,omitempty"`
	ProductID uint        `json:"productId"`
	Product   ProductInfo `json:"product,omitempty"`
	Comment   string      `json:"comment"`
	Rating    int         `json:"rating"`
	CreatedAt time.Time   `json:"createdAt"`
}

// UserInfo contains selected user fields for feedback response
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// ProductInfo contains selected product fields for feedback response
type ProductInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image"`
}
