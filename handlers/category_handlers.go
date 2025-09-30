package handlers

import (
	"net/http"
	"strconv"

	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
)

func CreateCategory(categoryService *service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := categoryService.CreateCategory(&category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func GetCategories(categoryService *service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := categoryService.GetAllCategories()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func GetCategory(categoryService *service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		category, err := categoryService.GetCategoryByID(uint(categoryID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusOK, category)
	}
}

func UpdateCategory(categoryService *service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		var req models.CategoryUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the request
		if err := models.ValidateStruct(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		// Update the category through service
		category, err := categoryService.UpdateCategory(uint(categoryID), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func DeleteCategory(categoryService *service.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		err = categoryService.DeleteCategory(uint(categoryID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category or category is not exist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	}
}
