package controllers

import (
	"net/http"
	"strconv"
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/services"

	"github.com/gin-gonic/gin"
)

type RecipeController struct {
	recipeService *services.RecipeService
}

func NewRecipeController(recipeService *services.RecipeService) *RecipeController {
	return &RecipeController{recipeService: recipeService}
}

type GenerateRecipesRequest struct {
	Cuisine string `json:"cuisine"` // optional cuisine preference
}

// GenerateRecipes generates new recipes based on user's groceries
func (c *RecipeController) GenerateRecipes(ctx *gin.Context) {
	userID := ctx.GetUint("userID") // Assuming you have auth middleware

	var req GenerateRecipesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	recipes, err := c.recipeService.GenerateRecipes(userID, req.Cuisine)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"recipes": recipes})
}

// GetAllRecipes returns all recipes for the authenticated user
func (c *RecipeController) GetAllRecipes(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	recipes, err := c.recipeService.GetAllRecipes(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"recipes": recipes})
}

// GetRecipeByID returns a specific recipe by ID for the authenticated user
func (c *RecipeController) GetRecipeByID(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	recipeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	recipe, err := c.recipeService.GetRecipeByID(userID, uint(recipeID))
	if err != nil {
		if err == models.ErrRecipeNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"recipe": recipe})
}
