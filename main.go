package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"health-store/handlers"
	"health-store/middleware"
	"health-store/models"
)

var DB *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "health-store.db"
	}

	DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection successful.")

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

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register(DB))
		authRoutes.POST("/login", handlers.Login(DB))
	}

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(DB, "admin"))
	{
		adminRoutes.GET("/users", handlers.GetUsers(DB))
		adminRoutes.GET("/users/:id", handlers.GetUser(DB))
		adminRoutes.PUT("/users/:id", handlers.UpdateUser(DB))
		adminRoutes.DELETE("/users/:id", handlers.DeleteUser(DB))

		adminRoutes.POST("/products", handlers.CreateProduct(DB))
		adminRoutes.GET("/products", handlers.GetProducts(DB))
		adminRoutes.GET("/products/:id", handlers.GetProduct(DB))
		adminRoutes.PUT("/products/:id", handlers.UpdateProduct(DB))
		adminRoutes.DELETE("/products/:id", handlers.DeleteProduct(DB))

		adminRoutes.POST("/categories", handlers.CreateCategory(DB))
		adminRoutes.GET("/categories", handlers.GetCategories(DB))
		adminRoutes.GET("/categories/:id", handlers.GetCategory(DB))
		adminRoutes.PUT("/categories/:id", handlers.UpdateCategory(DB))
		adminRoutes.DELETE("/categories/:id", handlers.DeleteCategory(DB))

		adminRoutes.GET("/report", handlers.GenerateReport)
	}

	cartRoutes := r.Group("/cart")
	cartRoutes.Use(middleware.AuthMiddleware(DB, "customer", "admin"))
	{
		cartRoutes.GET("/", handlers.GetCart(DB))
		cartRoutes.POST("/", handlers.AddToCart(DB))
		cartRoutes.DELETE("/:id", handlers.RemoveFromCart(DB))
	}

	orderRoutes := r.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware(DB, "customer", "admin"))
	{
		orderRoutes.POST("/", handlers.PlaceOrder(DB))
		orderRoutes.PUT("/:id/cancel", handlers.CancelOrder(DB))
	}

	feedbackRoutes := r.Group("/feedback")
	feedbackRoutes.Use(middleware.AuthMiddleware(DB, "customer", "admin"))
	{
		feedbackRoutes.POST("/", handlers.GiveFeedback(DB))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	fmt.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}