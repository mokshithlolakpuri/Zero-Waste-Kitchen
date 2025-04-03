package controllers

import (
	"net/http"
	"time"
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/services"
	"zero-waste-kitchen/pkg/database"

	"github.com/gin-gonic/gin"
)

func GetAllGroceries(c *gin.Context) {
	userID := c.GetUint("userID")

	var groceries []models.GroceryItem
	if err := database.DB.Where("user_id = ?", userID).Find(&groceries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groceries"})
		return
	}

	c.JSON(http.StatusOK, groceries)
}

func CreateGrocery(c *gin.Context) {
	userID := c.GetUint("userID")

	var grocery models.GroceryItem
	if err := c.ShouldBindJSON(&grocery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grocery.UserID = userID

	if err := database.DB.Create(&grocery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create grocery item"})
		return
	}

	c.JSON(http.StatusCreated, grocery)
}

func GetGrocery(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var grocery models.GroceryItem
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&grocery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grocery item not found"})
		return
	}

	c.JSON(http.StatusOK, grocery)
}

func UpdateGrocery(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var grocery models.GroceryItem
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&grocery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grocery item not found"})
		return
	}

	if err := c.ShouldBindJSON(&grocery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&grocery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update grocery item"})
		return
	}

	c.JSON(http.StatusOK, grocery)
}

func DeleteGrocery(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var grocery models.GroceryItem
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&grocery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grocery item not found"})
		return
	}

	if err := database.DB.Delete(&grocery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete grocery item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Grocery item deleted successfully"})
}

func GetExpiringGroceries(c *gin.Context) {
	userID := c.GetUint("userID")

	var groceries []models.GroceryItem
	threshold := time.Now().Add(7 * 24 * time.Hour) // Items expiring in next 7 days

	if err := database.DB.Where("user_id = ? AND expiry_date <= ?", userID, threshold).Find(&groceries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expiring groceries"})
		return
	}

	c.JSON(http.StatusOK, groceries)
}

func CheckExpiringItems() {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return
	}

	for _, user := range users {
		var expiringItems []models.GroceryItem
		threshold := time.Now().Add(7 * 24 * time.Hour)

		if err := database.DB.Where("user_id = ? AND expiry_date <= ?", user.ID, threshold).Find(&expiringItems).Error; err != nil {
			continue
		}

		if len(expiringItems) > 0 {
			services.SendExpiryNotification(user, expiringItems)
		}
	}
}
