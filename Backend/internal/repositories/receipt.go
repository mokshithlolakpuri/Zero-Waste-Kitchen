package repositories

import (
	"zero-waste-kitchen/internal/models"

	"gorm.io/gorm"
)

type ReceiptRepository interface {
	Create(receipt *models.Receipt) error
	FindByID(id uint, userID uint) (*models.Receipt, error)
	FindAll(userID uint) ([]models.Receipt, error)
}

type receiptRepository struct {
	db *gorm.DB
}

func NewReceiptRepository(db *gorm.DB) ReceiptRepository {
	return &receiptRepository{db: db}
}

func (r *receiptRepository) Create(receipt *models.Receipt) error {
	return r.db.Create(receipt).Error
}

func (r *receiptRepository) FindByID(id uint, userID uint) (*models.Receipt, error) {
	var receipt models.Receipt
	err := r.db.Preload("Items").Where("id = ? AND user_id = ?", id, userID).First(&receipt).Error
	return &receipt, err
}

func (r *receiptRepository) FindAll(userID uint) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := r.db.Preload("Items").Where("user_id = ?", userID).Find(&receipts).Error
	return receipts, err
}
