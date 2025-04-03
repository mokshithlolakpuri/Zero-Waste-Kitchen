package repositories

import (
	"time"
	"zero-waste-kitchen/internal/models"

	"gorm.io/gorm"
)

type GroceryRepository interface {
	Create(grocery *models.GroceryItem) error
	FindByID(id uint, userID uint) (*models.GroceryItem, error)
	FindAll(userID uint) ([]models.GroceryItem, error)
	Update(grocery *models.GroceryItem) error
	Delete(id uint) error
	FindExpiring(userID uint, threshold time.Time) ([]models.GroceryItem, error)
}

type groceryRepository struct {
	db *gorm.DB
}

func NewGroceryRepository(db *gorm.DB) GroceryRepository {
	return &groceryRepository{db: db}
}

func (r *groceryRepository) Create(grocery *models.GroceryItem) error {
	return r.db.Create(grocery).Error
}

func (r *groceryRepository) FindByID(id uint, userID uint) (*models.GroceryItem, error) {
	var grocery models.GroceryItem
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&grocery).Error
	return &grocery, err
}

func (r *groceryRepository) FindAll(userID uint) ([]models.GroceryItem, error) {
	var groceries []models.GroceryItem
	err := r.db.Where("user_id = ?", userID).Find(&groceries).Error
	return groceries, err
}

func (r *groceryRepository) Update(grocery *models.GroceryItem) error {
	return r.db.Save(grocery).Error
}

func (r *groceryRepository) Delete(id uint) error {
	return r.db.Delete(&models.GroceryItem{}, id).Error
}

func (r *groceryRepository) FindExpiring(userID uint, threshold time.Time) ([]models.GroceryItem, error) {
	var groceries []models.GroceryItem
	err := r.db.Where(
		"user_id = ? AND expiry_date <= ? AND expiry_date > ?",
		userID,
		threshold,
		time.Now(),
	).Find(&groceries).Error
	return groceries, err
}
