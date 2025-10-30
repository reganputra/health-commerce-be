package service

import (
	"health-store/models"
	"health-store/repositories"
)

// FeedbackService handles business logic for feedback
type FeedbackService struct {
	feedbackRepo repositories.FeedbackRepositoryInterface
}

// NewFeedbackService creates a new feedback service
func NewFeedbackService(feedbackRepo repositories.FeedbackRepositoryInterface) *FeedbackService {
	return &FeedbackService{feedbackRepo: feedbackRepo}
}

// CreateFeedback creates a new feedback
func (s *FeedbackService) CreateFeedback(feedback *models.Feedback) error {
	return s.feedbackRepo.Create(feedback)
}

// GetFeedbackByID gets a feedback by ID
func (s *FeedbackService) GetFeedbackByID(id uint) (*models.Feedback, error) {
	return s.feedbackRepo.FindByID(id)
}

// GetFeedbackByProductID gets feedback by product ID
func (s *FeedbackService) GetFeedbackByProductID(productID uint) ([]models.Feedback, error) {
	return s.feedbackRepo.FindByProductID(productID)
}

// GetFeedbackByUserID gets feedback by user ID
func (s *FeedbackService) GetFeedbackByUserID(userID uint) ([]models.Feedback, error) {
	return s.feedbackRepo.FindByUserID(userID)
}

// GetAllFeedback gets all feedback
func (s *FeedbackService) GetAllFeedback() ([]models.Feedback, error) {
	return s.feedbackRepo.FindAll()
}
