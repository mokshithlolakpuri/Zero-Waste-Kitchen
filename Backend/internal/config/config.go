package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
	ServerPort string
}

var AppConfig Config

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file - using environment variables")
	}

	// Set configuration values
	AppConfig = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "zero_waste_kitchen"),
		DBPort:     getEnv("DB_PORT", "5432"),
		JWTSecret:  getEnv("JWT_SECRET", "your_jwt_secret_key"),
		ServerPort: getEnv("PORT", "8080"),
	}

	// Validate required configurations
	validateConfig()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}

	return value
}

func validateConfig() {
	required := []struct {
		value     string
		name      string
		minLength int
	}{
		{AppConfig.JWTSecret, "JWT_SECRET", 32},
		{AppConfig.DBUser, "DB_USER", 1},
		{AppConfig.DBName, "DB_NAME", 1},
	}

	for _, req := range required {
		if len(req.value) < req.minLength {
			log.Fatalf("Configuration error: %s must be at least %d characters long", req.name, req.minLength)
		}
	}

	// Validate DB port is a valid number
	if _, err := strconv.Atoi(AppConfig.DBPort); err != nil {
		log.Fatalf("Configuration error: DB_PORT must be a valid number")
	}

	// Validate server port is a valid number
	if _, err := strconv.Atoi(AppConfig.ServerPort); err != nil {
		log.Fatalf("Configuration error: PORT must be a valid number")
	}
}

// GetDBConnectionString returns the formatted PostgreSQL connection string
func (c *Config) GetDBConnectionString() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=disable"
}

// GetServerAddress returns the formatted server address with port
func (c *Config) GetServerAddress() string {
	return ":" + c.ServerPort
}
