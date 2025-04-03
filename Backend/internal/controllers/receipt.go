package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/pkg/database"

	"github.com/gin-gonic/gin"
)

// receipt_controller.go
func UploadReceipt(c *gin.Context) {
	// Initialize variables
	var filePath string

	// Check if file was uploaded
	fileHeader, err := c.FormFile("file")
	if err == nil {
		// File was provided, save it
		filePath = "uploads/" + fileHeader.Filename
		if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse receipt data
	receiptJSON := c.PostForm("receipt")
	if receiptJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing receipt data"})
		return
	}

	var receiptData struct {
		StoreName    string  `json:"storeName"`
		PurchaseDate string  `json:"purchaseDate"`
		TotalAmount  float64 `json:"totalAmount"`
		Items        []struct {
			Name            string  `json:"name"`
			Quantity        float64 `json:"quantity"`
			Unit            string  `json:"unit"`
			Price           float64 `json:"price"`
			ExpiryDate      string  `json:"expiryDate"`
			StorageLocation string  `json:"storageLocation"`
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(receiptJSON), &receiptData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt data: " + err.Error()})
		return
	}

	// Parse dates
	purchaseDate, err := time.Parse("2006-01-02", receiptData.PurchaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase date format"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create receipt
	receipt := models.Receipt{
		UserID:       userID.(uint),
		ImagePath:    filePath,
		StoreName:    receiptData.StoreName,
		PurchaseDate: purchaseDate,
		TotalAmount:  receiptData.TotalAmount,
	}

	// Save receipt to database
	if err := database.DB.Create(&receipt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create receipt"})
		return
	}

	// Create grocery items
	for _, itemData := range receiptData.Items {
		expiryDate, err := time.Parse(time.RFC3339, itemData.ExpiryDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry date format"})
			return
		}

		item := models.GroceryItem{
			UserID:          userID.(uint),
			ReceiptID:       &receipt.ID,
			Name:            itemData.Name,
			Quantity:        itemData.Quantity,
			Unit:            itemData.Unit,
			ExpiryDate:      expiryDate,
			StorageLocation: itemData.StorageLocation,
		}

		if err := database.DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create grocery item"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Receipt and items saved successfully",
		"receipt": receipt,
	})
}

func GetAllReceipts(c *gin.Context) {
	userID := c.GetUint("userID")

	var receipts []models.Receipt
	if err := database.DB.Preload("Items").Where("user_id = ?", userID).Find(&receipts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch receipts"})
		return
	}

	c.JSON(http.StatusOK, receipts)
}

func GetReceipt(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var receipt models.Receipt
	if err := database.DB.Preload("Items").Where("id = ? AND user_id = ?", id, userID).First(&receipt).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

// processReceipt is a mock implementation of receipt processing
func processReceipt(imagePath string) ([]models.GroceryItem, error) {
	// In a real implementation, this would call an OCR API
	// For now, we'll return mock data
	return []models.GroceryItem{
		{
			Name:            "Milk",
			Quantity:        1,
			Unit:            "L",
			StorageLocation: "refrigerator",
			ManufactureDate: time.Now().AddDate(0, -1, 0), // 1 month ago
			ExpiryDate:      time.Now().AddDate(0, 0, 7),  // 7 days from now
		},
		{
			Name:            "Rice",
			Quantity:        2,
			Unit:            "kg",
			StorageLocation: "dry_pantry",
			ManufactureDate: time.Now().AddDate(0, -3, 0), // 3 months ago
			ExpiryDate:      time.Now().AddDate(1, 0, 0),  // 1 year from now
		},
	}, nil
}
