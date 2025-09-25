package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"health-store/models"
)

func GetCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		var cart models.Cart
		if err := db.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			// If cart is not found, it might not be an error. Return an empty cart.
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusOK, gin.H{"cart_items": []models.CartItem{}})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
			return
		}
		c.JSON(http.StatusOK, cart)
	}
}

func AddToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		var cartItem models.CartItem
		if err := c.ShouldBindJSON(&cartItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var cart models.Cart
		if err := db.Where("user_id = ?", userID).FirstOrCreate(&cart, models.Cart{UserID: userID}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create cart"})
			return
		}

		cartItem.CartID = cart.ID
		if err := db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}

		c.JSON(http.StatusOK, cartItem)
	}
}

func RemoveFromCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		cartItemID := c.Param("id")

		var cartItem models.CartItem
		if err := db.First(&cartItem, cartItemID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}

		var cart models.Cart
		if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cart not found for user"})
			return
		}

		if cartItem.CartID != cart.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to remove this item"})
			return
		}

		if err := db.Delete(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
	}
}