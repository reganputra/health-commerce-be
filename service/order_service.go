package service

import (
	"errors"
	"fmt"
	"health-store/models"
	"health-store/repositories"
	"math/rand"
	"time"
)

// OrderService handles business logic for orders
type OrderService struct {
	orderRepo   repositories.OrderRepositoryInterface
	cartRepo    repositories.CartRepositoryInterface
	productRepo repositories.ProductRepositoryInterface
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo repositories.OrderRepositoryInterface,
	cartRepo repositories.CartRepositoryInterface,
	productRepo repositories.ProductRepositoryInterface,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// PlaceOrder places a new order from user's cart with simulated payment
func (s *OrderService) PlaceOrder(userID uint, req models.PlaceOrderRequest) (*models.Order, error) {
	// Get user's cart with items and products
	cart, err := s.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found or empty")
	}

	// Check if cart has items
	if len(cart.CartItems) == 0 {
		return nil, errors.New("cannot place order with empty cart")
	}

	// Calculate total price and validate stock
	var totalPrice float64
	var orderItems []models.OrderItem

	// Collect all product IDs for batch loading (fixes N+1 query problem)
	productIDs := make([]uint, len(cart.CartItems))
	for i, cartItem := range cart.CartItems {
		productIDs[i] = cartItem.ProductID
	}

	// Batch load all products in a single query
	products, err := s.productRepo.FindByIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load products: %v", err)
	}

	// Create product map for efficient lookup
	productMap := make(map[uint]*models.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	// Process cart items using batch-loaded products
	for _, cartItem := range cart.CartItems {
		product, exists := productMap[cartItem.ProductID]
		if !exists {
			return nil, fmt.Errorf("product not found: %d", cartItem.ProductID)
		}

		// Check stock availability
		if product.Stock < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s (available: %d, requested: %d)",
				product.Name, product.Stock, cartItem.Quantity)
		}

		// Calculate item total
		itemTotal := product.Price * float64(cartItem.Quantity)
		totalPrice += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
		})
	}

	// Simulate payment processing based on payment method
	orderStatus, err := s.simulatePayment(req.PaymentMethod)
	if err != nil {
		return nil, err
	}

	// Create order
	order := &models.Order{
		UserID:        userID,
		Status:        orderStatus,
		TotalPrice:    totalPrice,
		PaymentMethod: req.PaymentMethod,
		BankName:      req.BankName,
	}

	// Create order first
	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Create order items (stock already reduced when added to cart)
	for _, item := range orderItems {
		item.OrderID = order.ID

		// Create order item
		err = s.orderRepo.CreateOrderItem(&item)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %v", err)
		}

		// Stock is already reduced when item was added to cart
		// No need to reduce again during order placement
	}

	// Clear cart
	err = s.cartRepo.ClearCart(cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart: %v", err)
	}

	return order, nil
}

// GetOrderByID gets an order by ID
func (s *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

// GetAllOrders gets all orders
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

// GetOrdersByUserID gets orders by user ID (for customer order history)
func (s *OrderService) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

// CancelOrder cancels an order and restores stock
func (s *OrderService) CancelOrder(orderID uint, userID uint) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("unauthorized to cancel this order")
	}

	if order.Status == "shipped" || order.Status == "cancelled" {
		return errors.New("cannot cancel order in current status")
	}

	// Get order items to restore stock
	orderItems, err := s.orderRepo.FindOrderItemsByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order items: %v", err)
	}

	// Restore stock for each item (use negative reduction to add back)
	for _, item := range orderItems {
		err = s.productRepo.ReduceStock(item.ProductID, -item.Quantity) // Negative = restore stock
		if err != nil {
			return fmt.Errorf("failed to restore stock for product %d: %v", item.ProductID, err)
		}
	}

	order.Status = "cancelled"
	return s.orderRepo.UpdateStatus(order.ID, "cancelled")
}

// UpdateOrderStatus updates order status (admin only)
func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, status) {
		return errors.New("invalid status transition")
	}

	order.Status = status
	return s.orderRepo.UpdateStatus(order.ID, status)
}

// isValidStatusTransition validates order status transitions
func (s *OrderService) isValidStatusTransition(from, to string) bool {
	transitions := map[string][]string{
		"pending":   {"paid", "cancelled"},
		"paid":      {"shipped", "cancelled"},
		"shipped":   {}, // Final state
		"cancelled": {}, // Final state
	}

	validStatuses, exists := transitions[from]
	if !exists {
		return false
	}

	for _, status := range validStatuses {
		if status == to {
			return true
		}
	}

	return false
}

// simulatePayment simulates payment processing for demo purposes
func (s *OrderService) simulatePayment(paymentMethod string) (string, error) {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	switch paymentMethod {
	case "cod":
		// Cash on delivery - always successful
		return "pending", nil

	case "paypal":
		// Simulate PayPal processing (95% success rate)
		if rand.Float32() < 0.95 {
			return "paid", nil
		}
		return "", errors.New("payment failed: insufficient funds")

	case "debit":
		// Simulate debit card processing (90% success rate)
		if rand.Float32() < 0.90 {
			return "paid", nil
		}
		return "", errors.New("payment failed: card declined")

	case "cc":
		// Simulate credit card processing (92% success rate)
		if rand.Float32() < 0.92 {
			return "paid", nil
		}
		return "", errors.New("payment failed: credit limit exceeded")

	default:
		return "", errors.New("unsupported payment method")
	}
}

// GetOrderStatistics gets order statistics for reporting
func (s *OrderService) GetOrderStatistics() (int64, error) {
	return s.orderRepo.GetOrderStatistics()
}
