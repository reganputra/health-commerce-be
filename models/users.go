package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"unique" json:"username" validate:"required,min=3,max=50"`
	Password      string    `gorm:"column:password" json:"-" validate:"required,min=6"`
	Email         string    `gorm:"unique" json:"email" validate:"required,email"`
	Dob           string    `json:"dob" validate:"required"`
	Gender        string    `json:"gender" validate:"required,oneof=M F"`
	Address       string    `json:"address" validate:"required,min=10,max=255"`
	City          string    `json:"city" validate:"required,min=2,max=100"`
	ContactNumber string    `json:"contact_number" validate:"required,min=10,max=15"`
	PaypalID      string    `json:"paypal_id"`
	Role          string    `json:"role" validate:"required,oneof=admin customer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UserRegisterRequest represents the request payload for user registration
type UserRegisterRequest struct {
	Username      string `json:"username" validate:"required,min=3,max=50"`
	Password      string `json:"password" validate:"required,min=6"`
	Email         string `json:"email" validate:"required,email"`
	Dob           string `gorm:"column:dob;type:date" json:"dob" validate:"required"`
	Gender        string `json:"gender" validate:"required,oneof=M F"`
	Address       string `json:"address" validate:"required,min=10,max=255"`
	City          string `json:"city" validate:"required,min=2,max=100"`
	ContactNumber string `json:"contact_number" validate:"required,min=10,max=15"`
}

type UserUpdateRequest struct {
	Username      *string `json:"username,omitempty"`
	Email         *string `json:"email,omitempty"`
	Dob           *string `json:"dob,omitempty"`
	Gender        *string `json:"gender,omitempty"`
	Address       *string `json:"address,omitempty"`
	City          *string `json:"city,omitempty"`
	ContactNumber *string `json:"contact_number,omitempty"`
	Role          *string `json:"role,omitempty"`
}

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using the validator
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
