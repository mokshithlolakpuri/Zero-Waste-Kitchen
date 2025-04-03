package database

import (
	"fmt"
	"log"
	"time"
	"zero-waste-kitchen/internal/config"
	"zero-waste-kitchen/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	dsn := config.AppConfig.GetDBConnectionString()

	// Configure GORM to use PostgreSQL
	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Adjust log level as needed
		NowFunc: func() time.Time {
			return time.Now().UTC() // Use UTC for all timestamps
		},
	}

	// Establish connection with retry logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), dbConfig)
		if err == nil {
			break
		}

		log.Printf("Attempt %d: Failed to connect to database. Retrying in 5 seconds...", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
}

// AutoMigrate runs database migrations
func AutoMigrate() {
	// First create independent tables
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate users: %v", err)
	}

	err = DB.AutoMigrate(&models.Receipt{})
	if err != nil {
		log.Fatalf("Failed to migrate receipts: %v", err)
	}

	// Then create dependent tables
	err = DB.AutoMigrate(&models.GroceryItem{})
	if err != nil {
		log.Fatalf("Failed to migrate grocery_items: %v", err)
	}

	// Create indexes
	err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_grocery_items_user_expiry ON grocery_items(user_id, expiry_date)").Error
	if err != nil {
		log.Printf("Failed to create index: %v", err)
	}

	log.Println("Database migration completed successfully")
}

// HealthCheck verifies database connectivity
func HealthCheck() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	return sqlDB.Ping()
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	return sqlDB.Close()
}
