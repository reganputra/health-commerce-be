package service

import (
	"errors"
	"health-store/models"
	"health-store/repositories"
)

// ShopService handles business logic for shops and shop requests
type ShopService struct {
	shopRequestRepo *repositories.ShopRequestRepository
	shopRepo        *repositories.ShopRepository
}

// NewShopService creates a new shop service
func NewShopService(shopRequestRepo *repositories.ShopRequestRepository, shopRepo *repositories.ShopRepository) *ShopService {
	return &ShopService{
		shopRequestRepo: shopRequestRepo,
		shopRepo:        shopRepo,
	}
}

// CreateShopRequest creates a new shop creation request (no restrictions on multiple requests/shops)
func (s *ShopService) CreateShopRequest(userID uint, req *models.ShopRequestCreateRequest) (*models.ShopRequest, error) {
	shopRequest := &models.ShopRequest{
		UserID:      userID,
		ShopName:    req.ShopName,
		Description: req.Description,
		Status:      "pending",
	}

	err := s.shopRequestRepo.Create(shopRequest)
	if err != nil {
		return nil, err
	}

	return shopRequest, nil
}

// GetShopRequestByID gets a shop request by ID
func (s *ShopService) GetShopRequestByID(id uint) (*models.ShopRequest, error) {
	return s.shopRequestRepo.FindByID(id)
}

// GetAllShopRequests gets all shop requests
func (s *ShopService) GetAllShopRequests() ([]models.ShopRequest, error) {
	return s.shopRequestRepo.FindAll()
}

// GetShopRequestsByStatus gets shop requests by status
func (s *ShopService) GetShopRequestsByStatus(status string) ([]models.ShopRequest, error) {
	return s.shopRequestRepo.FindByStatus(status)
}

// GetUserShopRequests gets shop requests for a specific user
func (s *ShopService) GetUserShopRequests(userID uint) ([]models.ShopRequest, error) {
	return s.shopRequestRepo.FindByUserID(userID)
}

// ApproveShopRequest approves a shop request and creates a shop
func (s *ShopService) ApproveShopRequest(requestID uint) error {
	// Get the shop request
	shopRequest, err := s.shopRequestRepo.FindByID(requestID)
	if err != nil {
		return err
	}

	// Check if already processed
	if shopRequest.Status != "pending" {
		return errors.New("shop request has already been processed")
	}

	// Create the shop (no restriction on multiple shops per user)
	shop := &models.Shop{
		UserID:      shopRequest.UserID,
		ShopName:    shopRequest.ShopName,
		Description: shopRequest.Description,
		IsActive:    true,
	}

	err = s.shopRepo.Create(shop)
	if err != nil {
		return err
	}

	// Update shop request status
	shopRequest.Status = "approved"
	err = s.shopRequestRepo.Update(shopRequest)
	if err != nil {
		return err
	}

	return nil
}

// RejectShopRequest rejects a shop request
func (s *ShopService) RejectShopRequest(requestID uint, reason string) error {
	shopRequest, err := s.shopRequestRepo.FindByID(requestID)
	if err != nil {
		return err
	}

	// Check if already processed
	if shopRequest.Status != "pending" {
		return errors.New("shop request has already been processed")
	}

	shopRequest.Status = "rejected"
	shopRequest.RejectionReason = reason

	return s.shopRequestRepo.Update(shopRequest)
}

// GetAllShops gets all shops
func (s *ShopService) GetAllShops() ([]models.Shop, error) {
	return s.shopRepo.FindAll()
}

// GetShopByID gets a shop by ID
func (s *ShopService) GetShopByID(id uint) (*models.Shop, error) {
	return s.shopRepo.FindByID(id)
}

// GetShopByUserID gets a shop by user ID
func (s *ShopService) GetShopByUserID(userID uint) (*models.Shop, error) {
	return s.shopRepo.FindByUserID(userID)
}

// GetAllShopsByUserID gets all shops by user ID
func (s *ShopService) GetAllShopsByUserID(userID uint) ([]models.Shop, error) {
	return s.shopRepo.FindAllByUserID(userID)
}

// UpdateShop updates a shop
func (s *ShopService) UpdateShop(shopID uint, shopName, description string) (*models.Shop, error) {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return nil, err
	}

	shop.ShopName = shopName
	shop.Description = description

	err = s.shopRepo.Update(shop)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

// DeleteShop deletes a shop
func (s *ShopService) DeleteShop(shopID uint) error {
	// Check if shop exists
	_, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return err
	}

	return s.shopRepo.Delete(shopID)
}

// ActivateShop activates a shop
func (s *ShopService) ActivateShop(shopID uint) (*models.Shop, error) {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return nil, err
	}

	if shop.IsActive {
		return nil, errors.New("shop is already active")
	}

	shop.IsActive = true
	err = s.shopRepo.Update(shop)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

// DeactivateShop deactivates a shop
func (s *ShopService) DeactivateShop(shopID uint) (*models.Shop, error) {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return nil, err
	}

	if !shop.IsActive {
		return nil, errors.New("shop is already inactive")
	}

	shop.IsActive = false
	err = s.shopRepo.Update(shop)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

// GetActiveShops gets all active shops
func (s *ShopService) GetActiveShops() ([]models.Shop, error) {
	return s.shopRepo.FindActiveShops()
}

// GetInactiveShops gets all inactive shops
func (s *ShopService) GetInactiveShops() ([]models.Shop, error) {
	return s.shopRepo.FindInactiveShops()
}
