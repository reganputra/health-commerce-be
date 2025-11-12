package handlers

import (
	"net/http"
	"strconv"

	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
)

// CreateGuestBookEntry allows visitors/guests to create a guestbook entry
func CreateGuestBookEntry(guestBookService *service.GuestBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.GuestBookCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the request
		if err := models.ValidateStruct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		entry, err := guestBookService.CreateEntry(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create guestbook entry"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Guestbook entry created successfully",
			"entry":   entry,
		})
	}
}

// GetAllGuestBookEntries allows admin to view all guestbook entries
func GetAllGuestBookEntries(guestBookService *service.GuestBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := guestBookService.GetAllEntries()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch guestbook entries"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"entries": entries,
			"count":   len(entries),
		})
	}
}

// GetGuestBookEntry allows viewing a specific guestbook entry
func GetGuestBookEntry(guestBookService *service.GuestBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		entryID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry ID"})
			return
		}

		entry, err := guestBookService.GetEntryByID(uint(entryID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Guestbook entry not found"})
			return
		}

		c.JSON(http.StatusOK, entry)
	}
}

// DeleteGuestBookEntry allows admin to delete a guestbook entry
func DeleteGuestBookEntry(guestBookService *service.GuestBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		entryID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry ID"})
			return
		}

		err = guestBookService.DeleteEntry(uint(entryID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete guestbook entry"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Guestbook entry deleted successfully"})
	}
}
