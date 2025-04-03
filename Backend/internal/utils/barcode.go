package utils

import (
	"fmt"
	"time"
	"zero-waste-kitchen/internal/models"
)

func GetProductInfo(barcode string) (*models.GroceryItem, error) {
	// In a real implementation, this would call a barcode API
	// For now, we'll return mock data

	if barcode == "" {
		return nil, fmt.Errorf("barcode is required")
	}

	// Mock response
	return &models.GroceryItem{
		Name:            "Sample Product",
		Barcode:         barcode,
		BatchNumber:     "BATCH123",
		ManufactureDate: time.Now().AddDate(0, -6, 0), // 6 months ago
		ExpiryDate:      time.Now().AddDate(0, 6, 0),  // 6 months from now
	}, nil
}
