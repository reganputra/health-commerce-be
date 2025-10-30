package service

import (
	"errors"
	"health-store/models"
	"health-store/repositories"
)

// CartService handles business logic for carts
type CartService struct {
	cartRepo    repositories.CartRepositoryInterface
	productRepo repositories.ProductRepositoryInterface
}

// NewCartService creates a new cart service
func NewCartService(cartRepo repositories.CartRepositoryInterface, productRepo repositories.ProductRepositoryInterface) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetOrCreateCart gets or creates a cart for a user
func (s *CartService) GetOrCreateCart(userID uint) (*models.Cart, error) {
	return s.cartRepo.FindOrCreateCart(userID)
}

// GetCartByUserID gets a cart by user ID
func (s *CartService) GetCartByUserID(userID uint) (*models.Cart, error) {
	return s.cartRepo.FindCartByUserID(userID)
}

// In AddToCart service method
func (s *CartService) AddToCart(userID uint, cartItem models.CartItem) error {
	// Validate product exists and has stock
	product, err := s.productRepo.FindByID(cartItem.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < cartItem.Quantity {
		return errors.New("insufficient stock")
	}

	// Reduce stock immediately
	err = s.productRepo.ReduceStock(cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		return err
	}

	// Create cart item
	cart, err := s.cartRepo.FindOrCreateCart(userID)
	if err != nil {
		// Rollback stock if cart creation fails
		s.productRepo.UpdateStock(cartItem.ProductID, product.Stock)
		return err
	}

	cartItem.CartID = cart.ID
	return s.cartRepo.CreateCartItem(&cartItem)
}

// RemoveFromCart removes an item from the cart
func (s *CartService) RemoveFromCart(cartItemID uint, userID uint) error {
	// Get cart item
	item, err := s.cartRepo.FindCartItemByID(cartItemID)
	if err != nil {
		return errors.New("cart item not found")
	}

	// Verify ownership
	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	if item.CartID != cart.ID {
		return errors.New("unauthorized to remove this item")
	}
	// Restore stock (add back the quantity)
	err = s.productRepo.ReduceStock(item.ProductID, -item.Quantity) // Negative = increase stock
	if err != nil {
		return err
	}

	return s.cartRepo.DeleteCartItem(cartItemID)
}

// ClearCart clears all items from a cart
func (s *CartService) ClearCart(userID uint) error {
	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return err
	}

	return s.cartRepo.ClearCart(cart.ID)
}
