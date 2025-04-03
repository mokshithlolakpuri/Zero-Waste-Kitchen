package utils

import (
	"time"
	"zero-waste-kitchen/internal/models"
)

func ProcessReceipt(imagePath string) ([]models.GroceryItem, error) {
	// In a real implementation, this would call an OCR API
	// For now, we'll return mock data

	// Mock response
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
