package handlers

import (
	"fmt"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
			return
		}

		// Validate required fields
		if feedback.ProductID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product ID is required. Make sure to send 'productId' (not 'product_id') in your request body",
			})
			return
		}

		if feedback.Comment == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment is required"})
			return
		}

		if feedback.Rating < 1 || feedback.Rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
			return
		}

		feedback.UserID = userID

		err := feedbackService.CreateFeedback(&feedback)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to give feedback: " + err.Error()})
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

// GetProductFeedback gets all feedback for a specific product
func GetProductFeedback(feedbackService *service.FeedbackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("productId")
		if productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
			return
		}

		// Convert productID to uint
		var id uint
		if _, err := fmt.Sscanf(productID, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		feedbacks, err := feedbackService.GetFeedbackByProductID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"feedbacks": feedbacks,
			"count":     len(feedbacks),
		})
	}
}
