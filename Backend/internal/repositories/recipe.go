package repositories

import (
	"errors"
	"zero-waste-kitchen/internal/models"

	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")

type RecipeRepository interface {
	Save(recipe *models.Recipe) error
	FindByUserID(userID uint) ([]models.Recipe, error)
	FindByID(userID uint, recipeID uint) (*models.Recipe, error)
	FindAll() ([]models.Recipe, error)
	DeleteByID(id uint) error
}

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

// Save inserts or updates a recipe in the database
func (r *recipeRepository) Save(recipe *models.Recipe) error {
	if err := r.db.Save(recipe).Error; err != nil {
		return err
	}
	return nil
}

// FindByID retrieves a recipe by its ID and user ID
func (r *recipeRepository) FindByID(userID uint, recipeID uint) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := r.db.Where("id = ? AND user_id = ?", recipeID, userID).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &recipe, nil
}

// FindByUserID retrieves all recipes for a specific user
func (r *recipeRepository) FindByUserID(userID uint) ([]models.Recipe, error) {
	var recipes []models.Recipe
	if err := r.db.Where("user_id = ?", userID).Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

// FindAll retrieves all recipes from the database
func (r *recipeRepository) FindAll() ([]models.Recipe, error) {
	var recipes []models.Recipe
	if err := r.db.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

// DeleteByID deletes a recipe by its ID
func (r *recipeRepository) DeleteByID(id uint) error {
	if err := r.db.Delete(&models.Recipe{}, id).Error; err != nil {
		return err
	}
	return nil
}
