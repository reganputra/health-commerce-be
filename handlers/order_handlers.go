package handlers

import (
	"net/http"
	"strconv"

	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
)

func PlaceOrder(orderService *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		var req models.PlaceOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
			return
		}

		// Validate the request
		if err := models.ValidateStruct(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		order, err := orderService.PlaceOrder(userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": order})
	}
}

// UpdateOrderStatus allows admin to update order status
func UpdateOrderStatus(orderService *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderIDStr := c.Param("id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		var req models.OrderStatusUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
			return
		}

		// Validate the request
		if err := models.ValidateStruct(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		err = orderService.UpdateOrderStatus(uint(orderID), req.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
	}
}

// GetAllOrders allows admin to view all orders
func GetAllOrders(orderService *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := orderService.GetAllOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}

// GetOrder allows viewing a specific order
func GetOrder(orderService *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		orderIDStr := c.Param("id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		order, err := orderService.GetOrderByID(uint(orderID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Check if user owns the order or is admin
		userRole := c.MustGet("userRole").(string)
		if order.UserID != userID && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view this order"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

// GetUserOrders allows customers to view their order history
func GetUserOrders(orderService *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		orders, err := orderService.GetOrdersByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"orders": orders,
			"count":  len(orders),
		})
	}
}

func CancelOrder(orderService *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		orderIDStr := c.Param("id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		err = orderService.CancelOrder(uint(orderID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
	}
}
