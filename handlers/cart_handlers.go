package handlers

import (
	"net/http"
	"strconv"

	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
)

func GetCart(cartService *service.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		cart, err := cartService.GetCartByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
			return
		}

		// Calculate total
		total := 0.0
		for _, item := range cart.CartItems {
			total += item.Product.Price * float64(item.Quantity)
		}

		// Add total to response
		response := gin.H{
			"cart":  cart,
			"total": total,
		}

		c.JSON(http.StatusOK, response)
	}
}

func AddToCart(cartService *service.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		var cartItem models.CartItem
		if err := c.ShouldBindJSON(&cartItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := cartService.AddToCart(userID, cartItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
	}
}

func RemoveFromCart(cartService *service.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		cartItemIDStr := c.Param("id")
		cartItemID, err := strconv.Atoi(cartItemIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
			return
		}

		err = cartService.RemoveFromCart(uint(cartItemID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
	}
}
