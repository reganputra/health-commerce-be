package repositories

import (
	"health-store/models"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Cache entry for products
type productCacheEntry struct {
	products  []models.Product
	timestamp time.Time
}

// ProductRepository handles database operations for products
type ProductRepository struct {
	db    *gorm.DB
	cache map[string]productCacheEntry
	mutex sync.RWMutex
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db:    db,
		cache: make(map[string]productCacheEntry),
	}
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

// FindAllCached finds all products with caching for better performance
func (r *ProductRepository) FindAllCached() ([]models.Product, error) {
	cacheKey := "products:all"

	// Check cache first
	r.mutex.RLock()
	if entry, exists := r.cache[cacheKey]; exists {
		// Check if cache is still valid (5 minutes TTL)
		if time.Since(entry.timestamp) < 5*time.Minute {
			r.mutex.RUnlock()
			return entry.products, nil
		}
	}
	r.mutex.RUnlock()

	// Cache miss or expired, fetch from database
	products, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	// Update cache
	r.mutex.Lock()
	r.cache[cacheKey] = productCacheEntry{
		products:  products,
		timestamp: time.Now(),
	}
	r.mutex.Unlock()

	return products, nil
}

// InvalidateCache clears the product cache (call after product updates)
func (r *ProductRepository) InvalidateCache() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.cache = make(map[string]productCacheEntry)
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

// FindByIDs finds multiple products by their IDs in a single query (optimizes N+1 problem)
func (r *ProductRepository) FindByIDs(ids []uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("id IN ?", ids).Find(&products).Error
	return products, err
}

// FindByIDsMap finds multiple products by their IDs and returns a map for efficient lookup
func (r *ProductRepository) FindByIDsMap(ids []uint) (map[uint]*models.Product, error) {
	products, err := r.FindByIDs(ids)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uint]*models.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}
	return productMap, nil
}
