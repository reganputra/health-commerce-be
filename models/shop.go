package models

import "time"

// ShopRequest represents a shop creation request from a customer
type ShopRequest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ShopName    string    `json:"shop_name" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"required,min=10,max=500"`
	Status      string    `json:"status" gorm:"default:'pending'" validate:"oneof=pending approved rejected"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Shop represents an approved shop
type Shop struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"` // User who owns the shop (can have multiple shops)
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ShopName    string    `json:"shop_name" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"required,min=10,max=500"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShopRequestCreateRequest represents the request to create a shop request
type ShopRequestCreateRequest struct {
	ShopName    string `json:"shop_name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=10,max=500"`
}

// ShopRequestApprovalRequest represents the admin's action on a shop request
type ShopRequestApprovalRequest struct {
	Status          string `json:"status" validate:"required,oneof=approved rejected"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

// ShopRequestResponse represents a shop request with selected user fields
type ShopRequestResponse struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	ShopName        string    `json:"shop_name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ShopUpdateRequest represents the request to update a shop
type ShopUpdateRequest struct {
	ShopName    string `json:"shop_name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=10,max=500"`
}
