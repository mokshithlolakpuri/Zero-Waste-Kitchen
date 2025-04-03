package models

import "errors"

type Recipe struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	UserID       uint   `json:"user_id"`
	Title        string `json:"title"`
	Ingredients  string `json:"ingredients"`  // Each ingredient on a new line
	Instructions string `json:"instructions"` // Step-by-step instructions
	PrepTime     int    `json:"prep_time"`    // In minutes
	CookTime     int    `json:"cook_time"`    // In minutes
	Servings     int    `json:"servings"`
	Difficulty   string `json:"difficulty"`
	Cuisine      string `json:"cuisine"`
}

var ErrRecipeNotFound = errors.New("recipe not found")
