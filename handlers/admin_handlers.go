package handlers

import (
	"net/http"
	"strconv"

	"health-store/service"

	"github.com/gin-gonic/gin"
)

// Users
func GetUsers(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userService.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUser(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := userService.GetUserByID(uint(userID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Get existing user
		existingUser, err := userService.GetUserByID(uint(userID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Bind update data
		var updateData map[string]interface{}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only the fields that are provided
		if username, ok := updateData["username"].(string); ok {
			existingUser.Username = username
		}
		if email, ok := updateData["email"].(string); ok {
			existingUser.Email = email
		}
		if dob, ok := updateData["dob"].(string); ok {
			existingUser.Dob = dob
		}
		if gender, ok := updateData["gender"].(string); ok {
			existingUser.Gender = gender
		}
		if address, ok := updateData["address"].(string); ok {
			existingUser.Address = address
		}
		if city, ok := updateData["city"].(string); ok {
			existingUser.City = city
		}
		if contactNumber, ok := updateData["contact_number"].(string); ok {
			existingUser.ContactNumber = contactNumber
		}
		if role, ok := updateData["role"].(string); ok {
			existingUser.Role = role
		}

		err = userService.UpdateUser(existingUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusOK, existingUser)
	}
}

func DeleteUser(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = userService.DeleteUser(uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
