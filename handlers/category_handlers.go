package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"health-store/models"
	"gorm.io/gorm"
)

func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func GetCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category
		if err := db.Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func GetCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusOK, category)
	}
}

func UpdateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
			return
		}
		c.JSON(http.StatusOK, category)
	}
}

func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	}
}