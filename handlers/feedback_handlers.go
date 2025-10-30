package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"health-store/models"
)

func GiveFeedback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		var feedback models.Feedback
		if err := c.ShouldBindJSON(&feedback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		feedback.UserID = userID

		if err := db.Create(&feedback).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to give feedback"})
			return
		}

		c.JSON(http.StatusOK, feedback)
	}
}