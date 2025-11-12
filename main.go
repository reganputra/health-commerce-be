package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"health-store/config"
	"health-store/models"
	"health-store/repositories"
	"health-store/routes"
	"health-store/service"
	"health-store/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/unidoc/unipdf/v3/common/license"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	err := godotenv.Load()
	if err != nil {
		utils.Warn("Error loading .env file, using environment variables")
	}

	// Initialize UniDoc PDF License (optional, for removing watermarks)
	if licenseKey := os.Getenv("UNIDOC_LICENSE_KEY"); licenseKey != "" {
		err := license.SetMeteredKey(licenseKey)
		if err != nil {
			log.Printf("Warning: Failed to set UniDoc license key: %v", err)
		} else {
			log.Println("UniDoc license key successfully set")
		}
	} else {
		log.Println("Warning: UNIDOC_LICENSE_KEY not set. PDFs will have watermarks or may fail to generate.")
		log.Println("To remove watermarks, set UNIDOC_LICENSE_KEY environment variable or get a free trial at: https://unidoc.io")
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
		&models.ShopRequest{},
		&models.Shop{},
		&models.GuestBook{},
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
	shopRequestRepo := repositories.NewShopRequestRepository(DB)
	shopRepo := repositories.NewShopRepository(DB)
	guestBookRepo := repositories.NewGuestBookRepository(DB)

	// Initialize Cloudinary service
	cloudinaryService, err := service.NewCloudinaryService(cfg.Storage.CloudinaryURL)
	if err != nil {
		utils.LogError(err, "Failed to initialize Cloudinary service")
		log.Fatal("Failed to initialize Cloudinary service:", err)
	}
	utils.Info("Cloudinary service initialized successfully")

	// Initialize services
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	feedbackService := service.NewFeedbackService(feedbackRepo)
	reportService := service.NewReportService(orderRepo, productRepo, userRepo)
	shopService := service.NewShopService(shopRequestRepo, shopRepo)
	guestBookService := service.NewGuestBookService(guestBookRepo)

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup all routes
	routes.SetupRoutes(
		r,
		DB,
		userService,
		productService,
		categoryService,
		orderService,
		cartService,
		feedbackService,
		reportService,
		cloudinaryService,
		shopService,
		guestBookService,
	)

	fmt.Printf("Starting server on port %s...\n", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
