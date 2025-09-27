package models

import (
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductCreateRequest represents the request payload for creating a product
type ProductCreateRequest struct {
	CategoryID  uint    `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	ImageURL    string  `json:"image_url" validate:"required,url"`
}

// ProductUpdateRequest represents the request payload for updating a product
type ProductUpdateRequest struct {
	CategoryID  uint    `json:"category_id,omitempty"`
	Name        string  `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description string  `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Price       float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
	ImageURL    string  `json:"image_url,omitempty" validate:"omitempty,url"`
}
