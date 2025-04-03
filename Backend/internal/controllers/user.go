package controllers

import (
	"log"
	"net/http"
	"strings"

	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/pkg/database"
	"zero-waste-kitchen/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register new user
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trim whitespace and validate inputs
	user.Email = strings.TrimSpace(user.Email)
	log.Printf("Password before trimming: %s", user.Password)
	user.Password = strings.TrimSpace(user.Password)
	log.Printf("Password after trimming: %s", user.Password)
	user.Name = strings.TrimSpace(user.Name)

	// Validate required fields
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Validate password length
	if len(user.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	log.Printf("Registering user with email: %s", user.Email)

	// Hash the password before storing
	if err := user.HashPassword(); err != nil {
		log.Printf("Hashing error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Save user to DB
	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Database create error: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	log.Printf("User registered successfully. ID: %d", user.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// Login user
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trim whitespace from input
	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User not found for email: %s", input.Email)
		} else {
			log.Printf("Database error: %v", err)
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Printf("Checking password for user: %s", user.Email)

	// Compare stored hash with input password
	if !user.ComparePassword(input.Password) {
		log.Printf("Password comparison failed for user: %s", user.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Printf("Authentication successful for user: %s", user.Email)

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Don't return the password hash
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't return password
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

// UpdateUser allows users to update their profile
func UpdateUser(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields only if they are provided
	if input.Name != "" {
		user.Name = strings.TrimSpace(input.Name)
	}
	if input.Email != "" {
		user.Email = strings.TrimSpace(input.Email)
	}
	if input.Password != "" {
		user.Password = strings.TrimSpace(input.Password)
		if err := user.HashPassword(); err != nil {
			log.Printf("Password hashing error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
	}

	// Save updated user
	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("User update error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Don't return password
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"message": "User updated successfully",
	})
}

// RegisterFCMToken allows users to register their FCM token for push notifications
func RegisterFCMToken(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("fcm_token", strings.TrimSpace(input.Token)).Error; err != nil {
		log.Printf("Failed to update FCM token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update FCM token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "FCM token registered successfully"})
}
