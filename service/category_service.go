package service

import (
	"errors"
	"health-store/models"
	"health-store/repositories"
)

// CategoryService handles business logic for categories
type CategoryService struct {
	categoryRepo repositories.CategoryRepositoryInterface
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repositories.CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.categoryRepo.Create(category)
}

// GetCategoryByID gets a category by ID
func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

// GetAllCategories gets all categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(id uint, req models.CategoryUpdateRequest) (*models.Category, error) {
	// Get existing category
	existingCategory, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update only the fields that are provided
	existingCategory.Name = req.Name
	existingCategory.Description = req.Description

	// Save the updated category
	err = s.categoryRepo.Update(existingCategory)
	if err != nil {
		return nil, err
	}

	return existingCategory, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category is not exist or has been deleted")
	}
	return s.categoryRepo.Delete(id)
}
