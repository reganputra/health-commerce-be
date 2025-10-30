package repositories

import "health-store/models"

// UserRepositoryInterface defines methods for user repository
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	FindAll() ([]models.User, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	// Report-specific methods
	GetUserCount() (int64, error)
}

// ProductRepositoryInterface defines methods for product repository
type ProductRepositoryInterface interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindByIDs(ids []uint) ([]models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	FindAll() ([]models.Product, error)
	FindByCategory(categoryID uint) ([]models.Product, error)
	UpdateStock(productID uint, quantity int) error
	ReduceStock(productID uint, quantity int) error
	// Report-specific methods
	GetTopSellingProducts(limit int) ([]models.TopProduct, error)
	GetProductCount() (int64, error)
}

// OrderRepositoryInterface defines methods for order repository
type OrderRepositoryInterface interface {
	Create(order *models.Order) error
	FindByID(id uint) (*models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order *models.Order) error
	UpdateStatus(orderID uint, status string) error
	UpdateOrderFields(orderID uint, updates map[string]interface{}) error
	CreateOrderItem(item *models.OrderItem) error
	GetOrderStatistics() (int64, error)
	GetDB() interface{} // For transactions
	FindOrderItemsByOrderID(orderID uint) ([]models.OrderItem, error)
	// Report-specific methods
	GetRecentOrders(limit int) ([]models.Order, error)
	GetTotalRevenue() (float64, error)
	GetOrdersByStatus() (map[string]int64, error)
	GetOrdersByDateRange(startDate, endDate string) ([]models.Order, error)
	GetTopCustomers(limit int) ([]models.TopCustomer, error)
	GetRevenueByDateRange(startDate, endDate string) (float64, error)
}

// CartRepositoryInterface defines methods for cart repository
type CartRepositoryInterface interface {
	FindOrCreateCart(userID uint) (*models.Cart, error)
	FindCartByUserID(userID uint) (*models.Cart, error)
	FindCartBasic(userID uint) (*models.Cart, error)
	FindCartWithCount(userID uint) (*models.Cart, int64, error)
	CreateCartItem(item *models.CartItem) error
	FindCartItemByID(id uint) (*models.CartItem, error)
	DeleteCartItem(id uint) error
	ClearCart(cartID uint) error
	GetCartItemCount(cartID uint) (int64, error)
}

// CategoryRepositoryInterface defines methods for category repository
type CategoryRepositoryInterface interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	FindAll() ([]models.Category, error)
}

// FeedbackRepositoryInterface defines methods for feedback repository
type FeedbackRepositoryInterface interface {
	Create(feedback *models.Feedback) error
	FindByID(id uint) (*models.Feedback, error)
	FindByProductID(productID uint) ([]models.Feedback, error)
	FindByUserID(userID uint) ([]models.Feedback, error)
	FindAll() ([]models.Feedback, error)
}
