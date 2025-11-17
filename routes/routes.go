package routes

import (
	"health-store/handlers"
	"health-store/middleware"
	"health-store/models"
	"health-store/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	r *gin.Engine,
	db *gorm.DB,
	userService *service.UserService,
	productService *service.ProductService,
	categoryService *service.CategoryService,
	orderService *service.OrderService,
	cartService *service.CartService,
	feedbackService *service.FeedbackService,
	reportService *service.ReportService,
	cloudinaryService *service.CloudinaryService,
	shopService *service.ShopService,
	guestBookService *service.GuestBookService,
) {
	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Setup route groups
	setupPublicRoutes(r, productService, categoryService, feedbackService)
	setupAuthRoutes(r, userService)
	setupAdminRoutes(r, db, userService, productService, categoryService, reportService, cloudinaryService, shopService, guestBookService, feedbackService)
	setupCartRoutes(r, db, cartService)
	setupOrderRoutes(r, db, orderService)
	setupAdminOrderRoutes(r, db, orderService)
	setupFeedbackRoutes(r, db, feedbackService)
	setupShopRoutes(r, db, shopService)
	setupGuestBookRoutes(r, guestBookService)

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
}

// setupPublicRoutes configures public routes accessible without authentication
func setupPublicRoutes(r *gin.Engine, productService *service.ProductService, categoryService *service.CategoryService, feedbackService *service.FeedbackService) {
	publicRoutes := r.Group("/api")
	{
		publicRoutes.GET("/products", handlers.GetProducts(productService))
		publicRoutes.GET("/products/:id", handlers.GetProduct(productService, feedbackService))
		publicRoutes.GET("/categories", handlers.GetCategories(categoryService))
		publicRoutes.GET("/categories/:id", handlers.GetCategory(categoryService))
	}
}

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(r *gin.Engine, userService *service.UserService) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register(userService))
		authRoutes.POST("/login", handlers.Login(userService))
	}
}

// setupAdminRoutes configures admin-only routes
func setupAdminRoutes(
	r *gin.Engine,
	db *gorm.DB,
	userService *service.UserService,
	productService *service.ProductService,
	categoryService *service.CategoryService,
	reportService *service.ReportService,
	cloudinaryService *service.CloudinaryService,
	shopService *service.ShopService,
	guestBookService *service.GuestBookService,
	feedbackService *service.FeedbackService,
) {
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(db, "admin"))
	{
		// User management
		adminRoutes.GET("/users", middleware.RequirePermission(models.PermissionReadUser), handlers.GetUsers(userService))
		adminRoutes.GET("/users/:id", middleware.RequirePermission(models.PermissionReadUser), handlers.GetUser(userService))
		adminRoutes.PUT("/users/:id", middleware.RequirePermission(models.PermissionUpdateUser), handlers.UpdateUser(userService))
		adminRoutes.DELETE("/users/:id", middleware.RequirePermission(models.PermissionDeleteUser), handlers.DeleteUser(userService))

		// Product management
		adminRoutes.POST("/products", middleware.RequirePermission(models.PermissionCreateProduct), handlers.CreateProduct(productService, cloudinaryService))
		adminRoutes.GET("/products", middleware.RequirePermission(models.PermissionReadProduct), handlers.GetProducts(productService))
		adminRoutes.GET("/products/:id", middleware.RequirePermission(models.PermissionReadProduct), handlers.GetProduct(productService, feedbackService))
		adminRoutes.PUT("/products/:id", middleware.RequirePermission(models.PermissionUpdateProduct), handlers.UpdateProduct(productService, cloudinaryService))
		adminRoutes.DELETE("/products/:id", middleware.RequirePermission(models.PermissionDeleteProduct), handlers.DeleteProduct(productService, cloudinaryService))

		// Category management
		adminRoutes.POST("/categories", middleware.RequirePermission(models.PermissionCreateCategory), handlers.CreateCategory(categoryService))
		adminRoutes.GET("/categories", middleware.RequirePermission(models.PermissionReadCategory), handlers.GetCategories(categoryService))
		adminRoutes.GET("/categories/:id", middleware.RequirePermission(models.PermissionReadCategory), handlers.GetCategory(categoryService))
		adminRoutes.PUT("/categories/:id", middleware.RequirePermission(models.PermissionUpdateCategory), handlers.UpdateCategory(categoryService))
		adminRoutes.DELETE("/categories/:id", middleware.RequirePermission(models.PermissionDeleteCategory), handlers.DeleteCategory(categoryService))

		// Reports
		adminRoutes.GET("/report", middleware.RequirePermission(models.PermissionReadReport), handlers.GenerateReport(reportService))

		// Shop request management
		adminRoutes.POST("/shop-requests", middleware.RequirePermission(models.PermissionCreateShopRequest), handlers.CreateShopRequest(shopService))
		adminRoutes.GET("/shop-requests", middleware.RequirePermission(models.PermissionReadShopRequest), handlers.GetAllShopRequests(shopService))
		adminRoutes.GET("/shop-requests/:id", middleware.RequirePermission(models.PermissionReadShopRequest), handlers.GetShopRequest(shopService))
		adminRoutes.PUT("/shop-requests/:id/approve", middleware.RequirePermission(models.PermissionApproveShop), handlers.ApproveShopRequest(shopService))
		adminRoutes.PUT("/shop-requests/:id/reject", middleware.RequirePermission(models.PermissionRejectShop), handlers.RejectShopRequest(shopService))

		// GuestBook management
		adminRoutes.GET("/guestbook", middleware.RequirePermission(models.PermissionReadGuestBook), handlers.GetAllGuestBookEntries(guestBookService))
		adminRoutes.GET("/guestbook/:id", middleware.RequirePermission(models.PermissionReadGuestBook), handlers.GetGuestBookEntry(guestBookService))
		adminRoutes.DELETE("/guestbook/:id", middleware.RequirePermission(models.PermissionDeleteGuestBook), handlers.DeleteGuestBookEntry(guestBookService))
	}
}

// setupCartRoutes configures cart routes
func setupCartRoutes(r *gin.Engine, db *gorm.DB, cartService *service.CartService) {
	cartRoutes := r.Group("/cart")
	cartRoutes.Use(middleware.AuthMiddleware(db, "customer", "admin"))
	cartRoutes.Use(middleware.RequirePermission(models.PermissionReadCart))
	{
		cartRoutes.GET("/", handlers.GetCart(cartService))
		cartRoutes.POST("/", middleware.RequirePermission(models.PermissionUpdateCart), handlers.AddToCart(cartService))
		cartRoutes.DELETE("/:id", middleware.RequirePermission(models.PermissionUpdateCart), handlers.RemoveFromCart(cartService))
	}
}

// setupOrderRoutes configures customer order routes
func setupOrderRoutes(r *gin.Engine, db *gorm.DB, orderService *service.OrderService) {
	orderRoutes := r.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware(db, "customer", "admin"))
	{
		orderRoutes.POST("/", middleware.RequirePermission(models.PermissionCreateOrder), handlers.PlaceOrder(orderService))
		orderRoutes.GET("/", handlers.GetUserOrders(orderService)) // Customer order history
		orderRoutes.GET("/:id", middleware.RequirePermission(models.PermissionReadOrder), handlers.GetOrder(orderService))
		orderRoutes.GET("/:id/receipt", middleware.RequirePermission(models.PermissionReadOrder), handlers.GeneratePurchaseReceipt(orderService))
		orderRoutes.PUT("/:id/cancel", middleware.RequirePermission(models.PermissionUpdateOrder), handlers.CancelOrder(orderService))
	}
}

// setupAdminOrderRoutes configures admin order management routes
func setupAdminOrderRoutes(r *gin.Engine, db *gorm.DB, orderService *service.OrderService) {
	adminOrderRoutes := r.Group("/admin/orders")
	adminOrderRoutes.Use(middleware.AuthMiddleware(db, "admin"))
	adminOrderRoutes.Use(middleware.RequirePermission(models.PermissionReadOrder))
	{
		adminOrderRoutes.GET("/", handlers.GetAllOrders(orderService))
		adminOrderRoutes.PUT("/:id/status", middleware.RequirePermission(models.PermissionUpdateOrder), handlers.UpdateOrderStatus(orderService))
	}
}

// setupFeedbackRoutes configures feedback routes
func setupFeedbackRoutes(r *gin.Engine, db *gorm.DB, feedbackService *service.FeedbackService) {
	feedbackRoutes := r.Group("/feedback")
	{
		// Public route to get feedback for a product
		feedbackRoutes.GET("/product/:productId", handlers.GetProductFeedback(feedbackService))

		// Protected route to submit feedback
		feedbackRoutes.POST("/", middleware.AuthMiddleware(db, "customer", "admin"), middleware.RequirePermission(models.PermissionCreateFeedback), handlers.GiveFeedback(feedbackService))
	}
}

// setupShopRoutes configures shop routes
func setupShopRoutes(r *gin.Engine, db *gorm.DB, shopService *service.ShopService) {
	shopRoutes := r.Group("/shops")
	{
		// Public routes
		shopRoutes.GET("/", handlers.GetAllShops(shopService))
		shopRoutes.GET("/active", handlers.GetActiveShops(shopService))
		shopRoutes.GET("/:id", handlers.GetShop(shopService))

		// Protected routes for authenticated users
		authenticated := shopRoutes.Group("")
		authenticated.Use(middleware.AuthMiddleware(db, "customer", "admin"))
		{
			// Get user's own shops
			authenticated.GET("/my/shops", handlers.GetMyShops(shopService))
			authenticated.GET("/my/shop", handlers.GetMyShop(shopService))

			// Update and delete (shop owners and admins)
			authenticated.PUT("/:id", handlers.UpdateShop(shopService))
			authenticated.DELETE("/:id", handlers.DeleteShop(shopService))
		}

		// Admin only routes
		adminShops := shopRoutes.Group("")
		adminShops.Use(middleware.AuthMiddleware(db, "admin"))
		{
			adminShops.GET("/inactive", middleware.RequirePermission(models.PermissionReadShop), handlers.GetInactiveShops(shopService))
			adminShops.PUT("/:id/activate", middleware.RequirePermission(models.PermissionUpdateShop), handlers.ActivateShop(shopService))
			adminShops.PUT("/:id/deactivate", middleware.RequirePermission(models.PermissionUpdateShop), handlers.DeactivateShop(shopService))
		}
	}
}

// setupGuestBookRoutes configures guest book routes
func setupGuestBookRoutes(r *gin.Engine, guestBookService *service.GuestBookService) {
	guestBookRoutes := r.Group("/guestbook")
	{
		// Public route for visitors to create entries
		guestBookRoutes.POST("/", handlers.CreateGuestBookEntry(guestBookService))
	}
}
