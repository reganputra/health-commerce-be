package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"health-store/models"
)

func PlaceOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		var cart models.Cart
		if err := db.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found or is empty"})
			return
		}

		if len(cart.CartItems) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot place an order with an empty cart"})
			return
		}

		var totalPrice float64
		var orderItems []models.OrderItem
		for _, item := range cart.CartItems {
			totalPrice += item.Product.Price * float64(item.Quantity)
			orderItems = append(orderItems, models.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
			})
		}

		order := models.Order{
			UserID:     userID,
			Status:     "pending",
			TotalPrice: totalPrice,
			OrderItems: orderItems,
		}

		// Use a transaction to ensure atomicity
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func CancelOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		orderID := c.Param("id")

		var order models.Order
		if err := db.First(&order, orderID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		if order.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to cancel this order"})
			return
		}

		if order.Status == "shipped" || order.Status == "cancelled" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot cancel an order that has been shipped or already cancelled"})
			return
		}

		order.Status = "cancelled"

		if err := db.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}