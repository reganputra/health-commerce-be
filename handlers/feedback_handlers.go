package handlers

import (
	"net/http"

	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
)

func GiveFeedback(feedbackService *service.FeedbackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		var feedback models.Feedback
		if err := c.ShouldBindJSON(&feedback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		feedback.UserID = userID

		err := feedbackService.CreateFeedback(&feedback)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to give feedback"})
			return
		}

		// Create response with selected fields
		response := models.FeedbackResponse{
			ID:        feedback.ID,
			UserID:    feedback.UserID,
			ProductID: feedback.ProductID,
			Comment:   feedback.Comment,
			Rating:    feedback.Rating,
			CreatedAt: feedback.CreatedAt,
		}

		if feedback.User.ID != 0 {
			response.User = models.UserInfo{
				ID:       feedback.User.ID,
				Username: feedback.User.Username,
			}
		}

		if feedback.Product.ID != 0 {
			response.Product = models.ProductInfo{
				ID:          feedback.Product.ID,
				Name:        feedback.Product.Name,
				Description: feedback.Product.Description,
				ImageURL:    feedback.Product.ImageURL,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}
