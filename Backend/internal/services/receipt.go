package services

import (
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/repositories"
)

type ReceiptService interface {
	CreateReceipt(receipt *models.Receipt) error
	GetReceiptByID(id uint, userID uint) (*models.Receipt, error)
	GetAllReceipts(userID uint) ([]models.Receipt, error)
}

type receiptService struct {
	repo repositories.ReceiptRepository
}

func NewReceiptService(repo repositories.ReceiptRepository) ReceiptService {
	return &receiptService{repo: repo}
}

func (s *receiptService) CreateReceipt(receipt *models.Receipt) error {
	return s.repo.Create(receipt)
}

func (s *receiptService) GetReceiptByID(id uint, userID uint) (*models.Receipt, error) {
	return s.repo.FindByID(id, userID)
}

func (s *receiptService) GetAllReceipts(userID uint) ([]models.Receipt, error) {
	return s.repo.FindAll(userID)
}
