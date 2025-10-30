package repositories

import (
	"health-store/models"

	"gorm.io/gorm"
)

// FeedbackRepository handles database operations for feedback
type FeedbackRepository struct {
	db *gorm.DB
}

// NewFeedbackRepository creates a new feedback repository
func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

// Create creates a new feedback
func (r *FeedbackRepository) Create(feedback *models.Feedback) error {
	return r.db.Create(feedback).Error
}

// FindByID finds a feedback by ID
func (r *FeedbackRepository) FindByID(id uint) (*models.Feedback, error) {
	var feedback models.Feedback
	err := r.db.Preload("User").Preload("Product").First(&feedback, id).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

// FindByProductID finds feedback by product ID
func (r *FeedbackRepository) FindByProductID(productID uint) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.db.Where("product_id = ?", productID).Preload("User").Find(&feedbacks).Error
	return feedbacks, err
}

// FindByUserID finds feedback by user ID
func (r *FeedbackRepository) FindByUserID(userID uint) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.db.Where("user_id = ?", userID).Preload("Product").Find(&feedbacks).Error
	return feedbacks, err
}

// FindAll finds all feedback
func (r *FeedbackRepository) FindAll() ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.db.Preload("User").Preload("Product").Find(&feedbacks).Error
	return feedbacks, err
}
