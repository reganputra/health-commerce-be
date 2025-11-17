package handlers

import (
	"net/http"
	"strconv"

	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
)

// CreateShopRequest allows an admin to create a shop creation request
func CreateShopRequest(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		var req models.ShopRequestCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the request
		if err := models.ValidateStruct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shopRequest, err := shopService.CreateShopRequest(userID, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Shop request submitted successfully",
			"request": shopRequest,
		})
	}
}

// GetMyShopRequests allows an admin to view their shop requests
func GetMyShopRequests(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		requests, err := shopService.GetUserShopRequests(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shop requests"})
			return
		}

		c.JSON(http.StatusOK, requests)
	}
}

// GetAllShopRequests allows admin to view all shop requests
func GetAllShopRequests(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Optional: filter by status
		status := c.Query("status")

		var requests []models.ShopRequest
		var err error

		if status != "" {
			requests, err = shopService.GetShopRequestsByStatus(status)
		} else {
			requests, err = shopService.GetAllShopRequests()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shop requests"})
			return
		}

		// Transform to response format
		responses := make([]models.ShopRequestResponse, len(requests))
		for i, req := range requests {
			responses[i] = models.ShopRequestResponse{
				ID:              req.ID,
				UserID:          req.UserID,
				Username:        req.User.Username,
				Email:           req.User.Email,
				ShopName:        req.ShopName,
				Description:     req.Description,
				Status:          req.Status,
				RejectionReason: req.RejectionReason,
				CreatedAt:       req.CreatedAt,
				UpdatedAt:       req.UpdatedAt,
			}
		}

		c.JSON(http.StatusOK, responses)
	}
}

// GetShopRequest allows admin to view a specific shop request
func GetShopRequest(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requestID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
			return
		}

		shopRequest, err := shopService.GetShopRequestByID(uint(requestID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop request not found"})
			return
		}

		response := models.ShopRequestResponse{
			ID:              shopRequest.ID,
			UserID:          shopRequest.UserID,
			Username:        shopRequest.User.Username,
			Email:           shopRequest.User.Email,
			ShopName:        shopRequest.ShopName,
			Description:     shopRequest.Description,
			Status:          shopRequest.Status,
			RejectionReason: shopRequest.RejectionReason,
			CreatedAt:       shopRequest.CreatedAt,
			UpdatedAt:       shopRequest.UpdatedAt,
		}

		c.JSON(http.StatusOK, response)
	}
}

// ApproveShopRequest allows admin to approve a shop request
func ApproveShopRequest(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requestID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
			return
		}

		err = shopService.ApproveShopRequest(uint(requestID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Shop request approved successfully"})
	}
}

// RejectShopRequest allows admin to reject a shop request
func RejectShopRequest(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requestID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
			return
		}

		var req models.ShopRequestApprovalRequest
		// Make the body optional - if not provided, use empty reason
		_ = c.ShouldBindJSON(&req)

		err = shopService.RejectShopRequest(uint(requestID), req.RejectionReason)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Shop request rejected successfully"})
	}
}

// GetAllShops allows viewing all shops
func GetAllShops(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shops, err := shopService.GetAllShops()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shops"})
			return
		}

		c.JSON(http.StatusOK, shops)
	}
}

// GetMyShop allows an admin to view their shop
func GetMyShop(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		shop, err := shopService.GetShopByUserID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
			return
		}

		c.JSON(http.StatusOK, shop)
	}
}

// GetMyShops allows a user to view all their shops
func GetMyShops(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		shops, err := shopService.GetAllShopsByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shops"})
			return
		}

		c.JSON(http.StatusOK, shops)
	}
}

// GetShop allows viewing a specific shop by ID
func GetShop(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		shopID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
			return
		}

		shop, err := shopService.GetShopByID(uint(shopID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop not found"})
			return
		}

		c.JSON(http.StatusOK, shop)
	}
}

// UpdateShop allows updating a shop
func UpdateShop(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		shopID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
			return
		}

		var req models.ShopUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the request
		if err := models.ValidateStruct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shop, err := shopService.UpdateShop(uint(shopID), req.ShopName, req.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Shop updated successfully",
			"shop":    shop,
		})
	}
}

// DeleteShop allows deleting a shop
func DeleteShop(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		shopID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
			return
		}

		err = shopService.DeleteShop(uint(shopID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Shop deleted successfully"})
	}
}

// ActivateShop allows activating a shop
func ActivateShop(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		shopID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
			return
		}

		shop, err := shopService.ActivateShop(uint(shopID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Shop activated successfully",
			"shop":    shop,
		})
	}
}

// DeactivateShop allows deactivating a shop
func DeactivateShop(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		shopID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID"})
			return
		}

		shop, err := shopService.DeactivateShop(uint(shopID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Shop deactivated successfully",
			"shop":    shop,
		})
	}
}

// GetActiveShops allows viewing all active shops
func GetActiveShops(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shops, err := shopService.GetActiveShops()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active shops"})
			return
		}

		c.JSON(http.StatusOK, shops)
	}
}

// GetInactiveShops allows viewing all inactive shops
func GetInactiveShops(shopService *service.ShopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shops, err := shopService.GetInactiveShops()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inactive shops"})
			return
		}

		c.JSON(http.StatusOK, shops)
	}
}
