package service

import (
	"health-store/models"
	"health-store/repositories"
)

// GuestBookService handles business logic for guest book
type GuestBookService struct {
	guestBookRepo *repositories.GuestBookRepository
}

// NewGuestBookService creates a new guest book service
func NewGuestBookService(guestBookRepo *repositories.GuestBookRepository) *GuestBookService {
	return &GuestBookService{guestBookRepo: guestBookRepo}
}

// CreateEntry creates a new guest book entry
func (s *GuestBookService) CreateEntry(req *models.GuestBookCreateRequest) (*models.GuestBook, error) {
	entry := &models.GuestBook{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	}

	err := s.guestBookRepo.Create(entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// GetAllEntries gets all guest book entries
func (s *GuestBookService) GetAllEntries() ([]models.GuestBook, error) {
	return s.guestBookRepo.FindAll()
}

// GetEntryByID gets a guest book entry by ID
func (s *GuestBookService) GetEntryByID(id uint) (*models.GuestBook, error) {
	return s.guestBookRepo.FindByID(id)
}

// DeleteEntry deletes a guest book entry
func (s *GuestBookService) DeleteEntry(id uint) error {
	return s.guestBookRepo.Delete(id)
}

// GetEntryCount gets the total count of guest book entries
func (s *GuestBookService) GetEntryCount() (int64, error) {
	return s.guestBookRepo.GetCount()
}
