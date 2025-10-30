package repositories

import (
	"health-store/models"

	"gorm.io/gorm"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindByID finds a product by ID
func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// FindAll finds all products
func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

// FindByCategory finds products by category ID
func (r *ProductRepository) FindByCategory(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// UpdateStock updates product stock
func (r *ProductRepository) UpdateStock(productID uint, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).
		Update("stock", quantity).Error
}

// ReduceStock reduces product stock by the specified quantity
func (r *ProductRepository) ReduceStock(productID uint, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error
}
