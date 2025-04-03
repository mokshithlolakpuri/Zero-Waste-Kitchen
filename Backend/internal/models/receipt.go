package models

import (
	"time"
)

type Receipt struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	UserID       uint          `json:"user_id"`
	ImagePath    string        `gorm:"not null" json:"image_path"`
	TotalAmount  float64       `json:"total_amount"`
	StoreName    string        `json:"store_name"`
	PurchaseDate time.Time     `json:"purchase_date"`
	CreatedAt    time.Time     `json:"created_at"`
	Items        []GroceryItem `gorm:"foreignKey:ReceiptID" json:"items"`
}
