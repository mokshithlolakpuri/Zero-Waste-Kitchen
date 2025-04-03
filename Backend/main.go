package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zero-waste-kitchen/internal/config"
	"zero-waste-kitchen/internal/controllers"
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/repositories"
	"zero-waste-kitchen/internal/services"
	"zero-waste-kitchen/pkg/database"
	"zero-waste-kitchen/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDB()
	database.AutoMigrate()
	defer database.CloseDB()

	// Assign initialized database instance to db
	db := database.DB

	// Initialize Firebase
	if err := services.InitializeFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize services
	groceryRepo := repositories.NewGroceryRepository(db)
	recipeRepo := repositories.NewRecipeRepository(db)
	recipeService := services.NewRecipeService(
		groceryRepo,
		recipeRepo,
		os.Getenv("GROQ_API_KEY"),
	)
	recipeController := controllers.NewRecipeController(recipeService)

	// Set Gin mode based on environment
	if config.AppConfig.ServerPort == "8080" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router with middleware
	router := gin.New()
	router.Use(
		middleware.CORSMiddleware(),
		middleware.RequestLogger(),
		gin.Recovery(),
	)

	// Register routes
	registerRoutes(router, recipeController)

	// Create HTTP server with graceful shutdown
	server := &http.Server{
		Addr:         config.AppConfig.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start notification service in background
	go startNotificationService()

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func registerRoutes(router *gin.Engine, recipeController *controllers.RecipeController) {
	api := router.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			if err := database.HealthCheck(); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.GET("/isadmin", middleware.JWTAuthMiddleware(), controllers.CheckAdminStatus)
		}

		// Admin routes
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.JWTAuthMiddleware(), middleware.RequireAdmin())
		{
			adminRoutes.GET("/users", controllers.GetUsersList)
			adminRoutes.POST("/send-notification", controllers.SendNotification)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			// Grocery routes
			grocery := protected.Group("/groceries")
			{
				grocery.GET("", controllers.GetAllGroceries)
				grocery.POST("", controllers.CreateGrocery)
				grocery.GET("/:id", controllers.GetGrocery)
				grocery.PUT("/:id", controllers.UpdateGrocery)
				grocery.DELETE("/:id", controllers.DeleteGrocery)
				grocery.GET("/expiring", controllers.GetExpiringGroceries)
			}

			// Receipt routes
			receipt := protected.Group("/receipts")
			{
				receipt.POST("/upload", controllers.UploadReceipt)
				receipt.GET("", controllers.GetAllReceipts)
				receipt.GET("/:id", controllers.GetReceipt)
			}

			// User routes
			user := protected.Group("/user")
			{
				user.GET("", controllers.GetCurrentUser)
				user.PUT("", controllers.UpdateUser)
				user.POST("/fcm-token", controllers.RegisterFCMToken)
			}

			// Recipe routes (new)
			recipe := protected.Group("/recipes")
			{
				recipe.GET("", recipeController.GetAllRecipes)
				recipe.GET("/:id", recipeController.GetRecipeByID)
				recipe.POST("/generate", recipeController.GenerateRecipes)
			}
		}
	}
}

func startNotificationService() {
	// Create tickers for different notification frequencies
	sevenDayTicker := time.NewTicker(24 * time.Hour)
	threeDayTicker := time.NewTicker(12 * time.Hour)
	oneDayTicker := time.NewTicker(4 * time.Hour)

	defer func() {
		sevenDayTicker.Stop()
		threeDayTicker.Stop()
		oneDayTicker.Stop()
	}()

	for {
		select {
		case <-sevenDayTicker.C:
			log.Println("Checking for items expiring in 7 days...")
			checkAndNotifyItems(7 * 24 * time.Hour)
		case <-threeDayTicker.C:
			log.Println("Checking for items expiring in 3 days...")
			checkAndNotifyItems(3 * 24 * time.Hour)
		case <-oneDayTicker.C:
			log.Println("Checking for items expiring in 1 day...")
			checkAndNotifyItems(24 * time.Hour)
		}
	}
}

func checkAndNotifyItems(threshold time.Duration) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Printf("Notification service error: %v", err)
		return
	}

	for _, user := range users {
		if user.FCMToken == "" {
			continue
		}

		var expiringItems []models.GroceryItem
		expiryThreshold := time.Now().Add(threshold)

		if err := database.DB.Where(
			"user_id = ? AND expiry_date <= ? AND expiry_date > ?",
			user.ID,
			expiryThreshold,
			time.Now(),
		).Find(&expiringItems).Error; err != nil {
			log.Printf("Failed to fetch expiring items for user %d: %v", user.ID, err)
			continue
		}

		if len(expiringItems) > 0 {
			services.SendExpiryNotification(user, expiringItems)
		}
	}
}
