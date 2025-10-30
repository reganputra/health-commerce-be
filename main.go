package main

import (
	"fmt"
	"log"

	"health-store/config"
	"health-store/handlers"
	"health-store/middleware"
	"health-store/models"
	"health-store/repositories"
	"health-store/service"
	"health-store/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	err := godotenv.Load()
	if err != nil {
		utils.Warn("Error loading .env file, using environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	utils.InitLogger(cfg)
	utils.Info("Starting Medical Equipment Online Store Backend")
	utils.Infof("Environment: %s", cfg.Server.Env)

	// Database connection
	DB, err = gorm.Open(mysql.Open(cfg.GetDatabaseDSN()), &gorm.Config{})
	if err != nil {
		utils.LogError(err, "Failed to connect to database")
		log.Fatal("Failed to connect to database:", err)
	}

	utils.Info("Database connection successful.")

	// Auto-migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Feedback{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migration successful.")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(DB)
	productRepo := repositories.NewProductRepository(DB)
	categoryRepo := repositories.NewCategoryRepository(DB)
	orderRepo := repositories.NewOrderRepository(DB)
	cartRepo := repositories.NewCartRepository(DB)
	feedbackRepo := repositories.NewFeedbackRepository(DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	feedbackService := service.NewFeedbackService(feedbackRepo)
	reportService := service.NewReportService(orderRepo, productRepo, userRepo)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Public routes for visitors
	publicRoutes := r.Group("/api")
	{
		publicRoutes.GET("/products", handlers.GetProducts(productService))
		publicRoutes.GET("/products/:id", handlers.GetProduct(productService))
		publicRoutes.GET("/categories", handlers.GetCategories(categoryService))
		publicRoutes.GET("/categories/:id", handlers.GetCategory(categoryService))
	}

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register(userService))
		authRoutes.POST("/login", handlers.Login(userService))
	}

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(DB, "admin"))
	{
		// User management
		adminRoutes.GET("/users", middleware.RequirePermission(models.PermissionReadUser), handlers.GetUsers(userService))
		adminRoutes.GET("/users/:id", middleware.RequirePermission(models.PermissionReadUser), handlers.GetUser(userService))
		adminRoutes.PUT("/users/:id", middleware.RequirePermission(models.PermissionUpdateUser), handlers.UpdateUser(userService))
		adminRoutes.DELETE("/users/:id", middleware.RequirePermission(models.PermissionDeleteUser), handlers.DeleteUser(userService))

		// Product management
		adminRoutes.POST("/products", middleware.RequirePermission(models.PermissionCreateProduct), handlers.CreateProduct(productService))
		adminRoutes.GET("/products", middleware.RequirePermission(models.PermissionReadProduct), handlers.GetProducts(productService))
		adminRoutes.GET("/products/:id", middleware.RequirePermission(models.PermissionReadProduct), handlers.GetProduct(productService))
		adminRoutes.PUT("/products/:id", middleware.RequirePermission(models.PermissionUpdateProduct), handlers.UpdateProduct(productService))
		adminRoutes.DELETE("/products/:id", middleware.RequirePermission(models.PermissionDeleteProduct), handlers.DeleteProduct(productService))

		// Category management
		adminRoutes.POST("/categories", middleware.RequirePermission(models.PermissionCreateCategory), handlers.CreateCategory(categoryService))
		adminRoutes.GET("/categories", middleware.RequirePermission(models.PermissionReadCategory), handlers.GetCategories(categoryService))
		adminRoutes.GET("/categories/:id", middleware.RequirePermission(models.PermissionReadCategory), handlers.GetCategory(categoryService))
		adminRoutes.PUT("/categories/:id", middleware.RequirePermission(models.PermissionUpdateCategory), handlers.UpdateCategory(categoryService))
		adminRoutes.DELETE("/categories/:id", middleware.RequirePermission(models.PermissionDeleteCategory), handlers.DeleteCategory(categoryService))

		// Reports
		adminRoutes.GET("/report", middleware.RequirePermission(models.PermissionReadReport), handlers.GenerateReport(reportService))
	}

	cartRoutes := r.Group("/cart")
	cartRoutes.Use(middleware.AuthMiddleware(DB, "customer", "admin"))
	cartRoutes.Use(middleware.RequirePermission(models.PermissionReadCart))
	{
		cartRoutes.GET("/", handlers.GetCart(cartService))
		cartRoutes.POST("/", middleware.RequirePermission(models.PermissionUpdateCart), handlers.AddToCart(cartService))
		cartRoutes.DELETE("/:id", middleware.RequirePermission(models.PermissionUpdateCart), handlers.RemoveFromCart(cartService))
	}

	orderRoutes := r.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware(DB, "customer", "admin"))
	{
		orderRoutes.POST("/", middleware.RequirePermission(models.PermissionCreateOrder), handlers.PlaceOrder(orderService))
		orderRoutes.GET("/", handlers.GetUserOrders(orderService)) // Customer order history
		orderRoutes.GET("/:id", middleware.RequirePermission(models.PermissionReadOrder), handlers.GetOrder(orderService))
		orderRoutes.PUT("/:id/cancel", middleware.RequirePermission(models.PermissionUpdateOrder), handlers.CancelOrder(orderService))
	}

	adminOrderRoutes := r.Group("/admin/orders")
	adminOrderRoutes.Use(middleware.AuthMiddleware(DB, "admin"))
	adminOrderRoutes.Use(middleware.RequirePermission(models.PermissionReadOrder))
	{
		adminOrderRoutes.GET("/", handlers.GetAllOrders(orderService))
		adminOrderRoutes.PUT("/:id/status", middleware.RequirePermission(models.PermissionUpdateOrder), handlers.UpdateOrderStatus(orderService))
	}

	feedbackRoutes := r.Group("/feedback")
	feedbackRoutes.Use(middleware.AuthMiddleware(DB, "customer", "admin"))
	feedbackRoutes.Use(middleware.RequirePermission(models.PermissionCreateFeedback))
	{
		feedbackRoutes.POST("/", handlers.GiveFeedback(feedbackService))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	fmt.Printf("Starting server on port %s...\n", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
