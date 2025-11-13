package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
)

func CreateProduct(productService *service.ProductService, cloudinaryService *service.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ProductCreateRequest

		// Check content type to determine binding method
		contentType := c.GetHeader("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// Manually parse form values
			if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
				return
			}

			// Debug: Log all form values
			categoryIDStr := c.PostForm("category_id")
			nameStr := c.PostForm("name")
			descStr := c.PostForm("description")
			priceStr := c.PostForm("price")
			stockStr := c.PostForm("stock")

			// Return debug info if description is empty
			if descStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Description is empty",
					"debug": gin.H{
						"category_id": categoryIDStr,
						"name": nameStr,
						"description": descStr,
						"price": priceStr,
						"stock": stockStr,
						"all_form_values": c.Request.PostForm,
					},
				})
				return
			}

			// Manually extract and convert form values
			categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
			price, _ := strconv.ParseFloat(priceStr, 64)
			stock, _ := strconv.Atoi(stockStr)

			req.CategoryID = uint(categoryID)
			req.Name = nameStr
			req.Description = descStr
			req.Price = price
			req.Stock = stock
			req.ImageURL = c.PostForm("image_url")

			// Handle image upload if provided
			file, header, err := c.Request.FormFile("image")
			if err == nil {
				defer file.Close()

				// Upload to Cloudinary
				imageURL, err := cloudinaryService.UploadImage(c.Request.Context(), file, header.Filename)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
					return
				}
				req.ImageURL = imageURL
			} else if err != http.ErrMissingFile {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process image: " + err.Error()})
				return
			}
		} else {
			// JSON binding for backward compatibility
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
				return
			}
		}

		// Validate the request
		if err := models.ValidateStruct(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		// Ensure we have an image URL
		if req.ImageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product image is required"})
			return
		}

		// Create product through service
		product, err := productService.CreateProduct(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func GetProducts(productService *service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := productService.GetAllProducts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func GetProduct(productService *service.ProductService, feedbackService *service.FeedbackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}
		product, err := productService.GetProductByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Get feedback for this product
		feedbacks, err := feedbackService.GetFeedbackByProductID(uint(id))
		if err != nil {
			// If error fetching feedback, return product without feedback
			c.JSON(http.StatusOK, product)
			return
		}

		// Create response with feedback
		response := gin.H{
			"id":          product.ID,
			"category_id": product.CategoryID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"image_url":   product.ImageURL,
			"created_at":  product.CreatedAt,
			"updated_at":  product.UpdatedAt,
			"feedbacks":   feedbacks,
		}

		// Include category if loaded
		if product.Category.ID != 0 {
			response["category"] = product.Category
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateProduct(productService *service.ProductService, cloudinaryService *service.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var req models.ProductUpdateRequest

		// Check content type to determine binding method
		contentType := c.GetHeader("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// Manually parse form values
			if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
				return
			}

			// Manually extract and convert form values (only if provided)
			if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
				categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
				req.CategoryID = uint(categoryID)
			}
			if name := c.PostForm("name"); name != "" {
				req.Name = name
			}
			if description := c.PostForm("description"); description != "" {
				req.Description = description
			}
			if priceStr := c.PostForm("price"); priceStr != "" {
				price, _ := strconv.ParseFloat(priceStr, 64)
				req.Price = price
			}
			if stockStr := c.PostForm("stock"); stockStr != "" {
				stock, _ := strconv.Atoi(stockStr)
				req.Stock = stock
			}
			if imageURL := c.PostForm("image_url"); imageURL != "" {
				req.ImageURL = imageURL
			}

			// Handle image upload if provided
			file, header, err := c.Request.FormFile("image")
			if err == nil {
				defer file.Close()

				// Get existing product to delete old image
				existingProduct, err := productService.GetProductByID(uint(id))
				if err == nil && existingProduct.ImageURL != "" {
					// Delete old image from Cloudinary
					_ = cloudinaryService.DeleteImage(c.Request.Context(), existingProduct.ImageURL)
				}

				// Upload new image to Cloudinary
				imageURL, err := cloudinaryService.UploadImage(c.Request.Context(), file, header.Filename)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
					return
				}
				req.ImageURL = imageURL
			} else if err != http.ErrMissingFile {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process image: " + err.Error()})
				return
			}
		} else {
			// JSON binding for backward compatibility
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
				return
			}
		}

		// Validate the request
		if err := models.ValidateStruct(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		// Update product through service
		product, err := productService.UpdateProduct(uint(id), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(productService *service.ProductService, cloudinaryService *service.CloudinaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		// Get product to delete image from Cloudinary
		product, err := productService.GetProductByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Delete product from database
		err = productService.DeleteProduct(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Delete image from Cloudinary (ignore errors as product is already deleted)
		if product.ImageURL != "" {
			_ = cloudinaryService.DeleteImage(c.Request.Context(), product.ImageURL)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}
