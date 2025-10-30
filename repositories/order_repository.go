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

// Update updates an order (updates all fields)
func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus updates only the status field for better performance
func (r *OrderRepository) UpdateStatus(orderID uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// UpdateOrderFields updates specific fields for better performance
func (r *OrderRepository) UpdateOrderFields(orderID uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(updates).Error
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

// GetRecentOrders returns the most recent orders with a limit
func (r *OrderRepository) GetRecentOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

// GetTotalRevenue calculates the total revenue from all orders
func (r *OrderRepository) GetTotalRevenue() (float64, error) {
	var totalRevenue float64
	err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(total_price), 0)").
		Where("status != ?", "cancelled").
		Scan(&totalRevenue).Error
	return totalRevenue, err
}

// GetOrdersByStatus returns count of orders grouped by status
func (r *OrderRepository) GetOrdersByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, result := range results {
		statusMap[result.Status] = result.Count
	}

	return statusMap, nil
}

// GetOrdersByDateRange returns orders within a date range
func (r *OrderRepository) GetOrdersByDateRange(startDate, endDate string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("User").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetTopCustomers returns top customers by order count and total spent
func (r *OrderRepository) GetTopCustomers(limit int) ([]models.TopCustomer, error) {
	var topCustomers []models.TopCustomer

	err := r.db.Model(&models.Order{}).
		Select("users.id as user_id, users.username, users.email, COUNT(orders.id) as order_count, COALESCE(SUM(orders.total_price), 0) as total_spent").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.status != ?", "cancelled").
		Group("users.id, users.username, users.email").
		Order("total_spent DESC").
		Limit(limit).
		Scan(&topCustomers).Error

	return topCustomers, err
}

// GetRevenueByDateRange calculates revenue within a date range
func (r *OrderRepository) GetRevenueByDateRange(startDate, endDate string) (float64, error) {
	var revenue float64
	err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(total_price), 0)").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status != ?", "cancelled").
		Scan(&revenue).Error
	return revenue, err
}
