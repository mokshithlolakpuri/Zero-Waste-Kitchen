package models

import (
	"time"
)

type StorageLocation string

const (
	DeepFreeze   StorageLocation = "deep_freeze"
	Refrigerator StorageLocation = "refrigerator"
	DryPantry    StorageLocation = "dry_pantry"
)

type GroceryItem struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `json:"user_id"`
	ReceiptID       *uint     `json:"receipt_id,omitempty"` // Make this a pointer to allow null values
	Name            string    `gorm:"not null" json:"name"`
	Quantity        float64   `gorm:"not null" json:"quantity"`
	Unit            string    `json:"unit"`
	Barcode         string    `json:"barcode"`
	BatchNumber     string    `json:"batch_number"`
	ManufactureDate time.Time `json:"manufacture_date"`
	ExpiryDate      time.Time `json:"expiry_date"`
	StorageLocation string    `gorm:"not null" json:"storageLocation"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
