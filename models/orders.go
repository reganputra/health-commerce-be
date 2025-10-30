package models

import (
	"time"
)

type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	UserID        uint        `gorm:"column:user_id;not null;index" json:"user_id"`
	User          User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Status        string      `gorm:"column:status;not null;index" json:"status"`
	TotalPrice    float64     `gorm:"column:total_price;not null" json:"total_price"`
	PaymentMethod string      `gorm:"column:payment_method;not null" json:"payment_method"`
	BankName      string      `gorm:"column:bank_name" json:"bank_name,omitempty"`
	CreatedAt     time.Time   `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
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
