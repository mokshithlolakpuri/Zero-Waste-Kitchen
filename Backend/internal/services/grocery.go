package services

import (
	"errors"
	"time"
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/repositories"
	"zero-waste-kitchen/internal/utils"
)

type GroceryService interface {
	CreateGrocery(grocery *models.GroceryItem) error
	GetGroceryByID(id uint, userID uint) (*models.GroceryItem, error)
	GetAllGroceries(userID uint) ([]models.GroceryItem, error)
	UpdateGrocery(grocery *models.GroceryItem, userID uint) error
	DeleteGrocery(id uint, userID uint) error
	GetExpiringGroceries(userID uint, days int) ([]models.GroceryItem, error)
}

type groceryService struct {
	repo repositories.GroceryRepository
}

func NewGroceryService(repo repositories.GroceryRepository) GroceryService {
	return &groceryService{repo: repo}
}

func (s *groceryService) CreateGrocery(grocery *models.GroceryItem) error {
	if err := utils.ValidateStruct(grocery); err != nil {
		return err
	}

	if grocery.ExpiryDate.Before(time.Now()) {
		return errors.New("expiry date cannot be in the past")
	}

	return s.repo.Create(grocery)
}

func (s *groceryService) GetGroceryByID(id uint, userID uint) (*models.GroceryItem, error) {
	return s.repo.FindByID(id, userID)
}

func (s *groceryService) GetAllGroceries(userID uint) ([]models.GroceryItem, error) {
	return s.repo.FindAll(userID)
}

func (s *groceryService) UpdateGrocery(grocery *models.GroceryItem, userID uint) error {
	existing, err := s.repo.FindByID(grocery.ID, userID)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return errors.New("unauthorized to update this item")
	}

	return s.repo.Update(grocery)
}

func (s *groceryService) DeleteGrocery(id uint, userID uint) error {
	existing, err := s.repo.FindByID(id, userID)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return errors.New("unauthorized to delete this item")
	}

	return s.repo.Delete(id)
}

func (s *groceryService) GetExpiringGroceries(userID uint, days int) ([]models.GroceryItem, error) {
	threshold := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	return s.repo.FindExpiring(userID, threshold)
}
