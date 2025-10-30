package repositories

import (
	"health-store/models"

	"gorm.io/gorm"
)

// CartRepository handles database operations for carts
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// FindOrCreateCart finds a user's cart or creates one if it doesn't exist
func (r *CartRepository) FindOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new cart
			cart = models.Cart{UserID: userID}
			if err := r.db.Create(&cart).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &cart, nil
}

// FindCartByUserID finds a cart with minimal data for frontend
// FindCartByUserID finds a cart by user ID
func (r *CartRepository) FindCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.
		// Preload("User").
		Preload("CartItems.Product").
		// Preload("CartItems.Product.Category").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// CreateCartItem creates a new cart item
func (r *CartRepository) CreateCartItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

// FindCartItemByID finds a cart item by ID
func (r *CartRepository) FindCartItemByID(id uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Preload("Product").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// DeleteCartItem deletes a cart item
func (r *CartRepository) DeleteCartItem(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}

// ClearCart removes all items from a cart
func (r *CartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

// GetCartItemCount returns the number of items in a cart
func (r *CartRepository) GetCartItemCount(cartID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.CartItem{}).Where("cart_id = ?", cartID).Count(&count).Error
	return count, err
}
