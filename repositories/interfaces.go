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
}

// ProductRepositoryInterface defines methods for product repository
type ProductRepositoryInterface interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	FindAll() ([]models.Product, error)
	FindByCategory(categoryID uint) ([]models.Product, error)
	UpdateStock(productID uint, quantity int) error
	ReduceStock(productID uint, quantity int) error
}

// OrderRepositoryInterface defines methods for order repository
type OrderRepositoryInterface interface {
	Create(order *models.Order) error
	FindByID(id uint) (*models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order *models.Order) error
	CreateOrderItem(item *models.OrderItem) error
	GetOrderStatistics() (int64, error)
}

// CartRepositoryInterface defines methods for cart repository
type CartRepositoryInterface interface {
	FindOrCreateCart(userID uint) (*models.Cart, error)
	FindCartByUserID(userID uint) (*models.Cart, error)
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
