package service

import (
	"errors"
	"health-store/models"
	"health-store/repositories"
)

// ProductService handles business logic for products
type ProductService struct {
	productRepo  repositories.ProductRepositoryInterface
	categoryRepo repositories.CategoryRepositoryInterface
}

// NewProductService creates a new product service
func NewProductService(productRepo repositories.ProductRepositoryInterface, categoryRepo repositories.CategoryRepositoryInterface) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req models.ProductCreateRequest) (*models.Product, error) {
	// Validate that category exists
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID gets a product by ID
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

// GetAllProducts gets all products
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.FindAll()
}

// GetProductsByCategory gets products by category
func (s *ProductService) GetProductsByCategory(categoryID uint) ([]models.Product, error) {
	return s.productRepo.FindByCategory(categoryID)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(id uint, req models.ProductUpdateRequest) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if category exists if CategoryID is provided
	if req.CategoryID != 0 {
		_, err = s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = req.CategoryID
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Stock != 0 || req.Stock == 0 { // Allow setting stock to 0
		product.Stock = req.Stock
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

// CheckProductStock checks if product has sufficient stock
func (s *ProductService) CheckProductStock(productID uint, quantity int) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return err
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	return nil
}

// ReduceProductStock reduces product stock
func (s *ProductService) ReduceProductStock(productID uint, quantity int) error {
	return s.productRepo.ReduceStock(productID, quantity)
}
