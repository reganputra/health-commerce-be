package repositories

import (
	"health-store/models"

	"gorm.io/gorm"
)

// GuestBookRepository handles database operations for guest book entries
type GuestBookRepository struct {
	db *gorm.DB
}

// NewGuestBookRepository creates a new guest book repository
func NewGuestBookRepository(db *gorm.DB) *GuestBookRepository {
	return &GuestBookRepository{db: db}
}

// Create creates a new guest book entry
func (r *GuestBookRepository) Create(entry *models.GuestBook) error {
	return r.db.Create(entry).Error
}

// FindByID finds a guest book entry by ID
func (r *GuestBookRepository) FindByID(id uint) (*models.GuestBook, error) {
	var entry models.GuestBook
	err := r.db.First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// FindAll finds all guest book entries
func (r *GuestBookRepository) FindAll() ([]models.GuestBook, error) {
	var entries []models.GuestBook
	err := r.db.Order("created_at DESC").Find(&entries).Error
	return entries, err
}

// Delete deletes a guest book entry
func (r *GuestBookRepository) Delete(id uint) error {
	return r.db.Delete(&models.GuestBook{}, id).Error
}

// GetCount returns the total count of guest book entries
func (r *GuestBookRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.GuestBook{}).Count(&count).Error
	return count, err
}
