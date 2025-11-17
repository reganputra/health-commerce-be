package repositories

import (
	"health-store/models"

	"gorm.io/gorm"
)

// ShopRequestRepository handles database operations for shop requests
type ShopRequestRepository struct {
	db *gorm.DB
}

// NewShopRequestRepository creates a new shop request repository
func NewShopRequestRepository(db *gorm.DB) *ShopRequestRepository {
	return &ShopRequestRepository{db: db}
}

// Create creates a new shop request
func (r *ShopRequestRepository) Create(request *models.ShopRequest) error {
	return r.db.Create(request).Error
}

// FindByID finds a shop request by ID
func (r *ShopRequestRepository) FindByID(id uint) (*models.ShopRequest, error) {
	var request models.ShopRequest
	err := r.db.Preload("User").First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// FindByUserID finds shop requests by user ID
func (r *ShopRequestRepository) FindByUserID(userID uint) ([]models.ShopRequest, error) {
	var requests []models.ShopRequest
	err := r.db.Where("user_id = ?", userID).Preload("User").Find(&requests).Error
	return requests, err
}

// FindAll finds all shop requests
func (r *ShopRequestRepository) FindAll() ([]models.ShopRequest, error) {
	var requests []models.ShopRequest
	err := r.db.Preload("User").Order("created_at DESC").Find(&requests).Error
	return requests, err
}

// FindByStatus finds shop requests by status
func (r *ShopRequestRepository) FindByStatus(status string) ([]models.ShopRequest, error) {
	var requests []models.ShopRequest
	err := r.db.Where("status = ?", status).Preload("User").Order("created_at DESC").Find(&requests).Error
	return requests, err
}

// Update updates a shop request
func (r *ShopRequestRepository) Update(request *models.ShopRequest) error {
	return r.db.Save(request).Error
}

// HasPendingRequest checks if a user has a pending shop request
func (r *ShopRequestRepository) HasPendingRequest(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ShopRequest{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&count).Error
	return count > 0, err
}

// ShopRepository handles database operations for shops
type ShopRepository struct {
	db *gorm.DB
}

// NewShopRepository creates a new shop repository
func NewShopRepository(db *gorm.DB) *ShopRepository {
	return &ShopRepository{db: db}
}

// Create creates a new shop
func (r *ShopRepository) Create(shop *models.Shop) error {
	return r.db.Create(shop).Error
}

// FindByID finds a shop by ID
func (r *ShopRepository) FindByID(id uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Preload("User").First(&shop, id).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

// FindByUserID finds a shop by user ID (returns first shop if user has multiple)
func (r *ShopRepository) FindByUserID(userID uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&shop).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

// FindAllByUserID finds all shops by user ID
func (r *ShopRepository) FindAllByUserID(userID uint) ([]models.Shop, error) {
	var shops []models.Shop
	err := r.db.Where("user_id = ?", userID).Preload("User").Find(&shops).Error
	return shops, err
}

// FindAll finds all shops
func (r *ShopRepository) FindAll() ([]models.Shop, error) {
	var shops []models.Shop
	err := r.db.Preload("User").Find(&shops).Error
	return shops, err
}

// Update updates a shop
func (r *ShopRepository) Update(shop *models.Shop) error {
	return r.db.Save(shop).Error
}

// Delete deletes a shop
func (r *ShopRepository) Delete(id uint) error {
	return r.db.Delete(&models.Shop{}, id).Error
}

// ExistsByUserID checks if a shop exists for a user
func (r *ShopRepository) ExistsByUserID(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Shop{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}

// FindByActiveStatus finds shops by active status
func (r *ShopRepository) FindByActiveStatus(isActive bool) ([]models.Shop, error) {
	var shops []models.Shop
	err := r.db.Where("is_active = ?", isActive).Preload("User").Find(&shops).Error
	return shops, err
}

// FindActiveShops finds all active shops
func (r *ShopRepository) FindActiveShops() ([]models.Shop, error) {
	return r.FindByActiveStatus(true)
}

// FindInactiveShops finds all inactive shops
func (r *ShopRepository) FindInactiveShops() ([]models.Shop, error) {
	return r.FindByActiveStatus(false)
}
