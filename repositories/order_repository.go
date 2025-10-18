package repositories

import (
	"health-store/models"

	"gorm.io/gorm"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// FindByID finds an order by ID
func (r *OrderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("OrderItems.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByUserID finds orders by user ID
func (r *OrderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("User").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Find(&orders).Error
	return orders, err
}

// FindAll finds all orders
func (r *OrderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("User").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Find(&orders).Error
	return orders, err
}

// Update updates an order
func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// CreateOrderItem creates an order item
func (r *OrderRepository) CreateOrderItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

// GetOrderStatistics returns order statistics for reporting
func (r *OrderRepository) GetOrderStatistics() (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Count(&count).Error
	return count, err
}

// GetDB returns the database instance for transactions
func (r *OrderRepository) GetDB() interface{} {
	return r.db
}

// FindOrderItemsByOrderID finds all order items for a specific order
func (r *OrderRepository) FindOrderItemsByOrderID(orderID uint) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	err := r.db.Where("order_id = ?", orderID).Preload("Product").Find(&orderItems).Error
	return orderItems, err
}
