package models

import "time"

// GuestBook represents a guest book entry from visitors
type GuestBook struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Message   string    `json:"message" validate:"required,min=10,max=1000"`
	CreatedAt time.Time `json:"created_at"`
}

// GuestBookCreateRequest represents the request to create a guestbook entry
type GuestBookCreateRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"required,email"`
	Message string `json:"message" validate:"required,min=10,max=1000"`
}
