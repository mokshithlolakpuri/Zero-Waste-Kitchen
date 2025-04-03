package controllers

import (
	"fmt"
	"log"
	"net/http"

	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/services"
	"zero-waste-kitchen/pkg/database"

	"github.com/gin-gonic/gin"
)

// CheckAdminStatus returns whether the user is an admin
func CheckAdminStatus(c *gin.Context) {
	// Get userID from the context (set by JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch the user from the database
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Return the is_admin status
	c.JSON(http.StatusOK, gin.H{"is_admin": user.IsAdmin})
}

// GetUsersList fetches all users
func GetUsersList(c *gin.Context) {
	// Define a struct to hold the filtered user data
	type UserResponse struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		FCMToken string `json:"fcm_token"`
		IsAdmin  bool   `json:"is_admin"`
	}

	var users []models.User
	if err := database.DB.Order("id ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Map the users to the filtered response
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			FCMToken: user.FCMToken,
			IsAdmin:  user.IsAdmin,
		})
	}

	// Return the filtered user data
	c.JSON(http.StatusOK, userResponses)
}

// SendNotification sends a push notification to a user
func SendNotification(c *gin.Context) {
	var input struct {
		UserID  uint   `json:"userId" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.FCMToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not have an FCM token"})
		return
	}
	log.Printf("notification's users FCM token: %s", user.FCMToken)

	if err := services.SendPushNotification(user.FCMToken, input.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send notification: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}
