package models

import (
	"time"
)

type Order struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status        string    `json:"status"`
	TotalPrice    float64   `json:"total_price"`
	PaymentMethod string    `json:"payment_method"`
	BankName      string    `json:"bank_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PlaceOrderRequest represents the request payload for placing an order
type PlaceOrderRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=paypal debit cc cod"`
	BankName      string `json:"bank_name,omitempty"`
}

// OrderStatusUpdateRequest represents the request payload for updating order status (admin only)
type OrderStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=pending paid shipped cancelled"`
}
