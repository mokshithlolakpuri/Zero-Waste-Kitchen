package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/repositories"

	"github.com/jpoz/groq"
)

type RecipeService struct {
	groceryRepo repositories.GroceryRepository
	recipeRepo  repositories.RecipeRepository
	groqClient  *groq.Client
}

func NewRecipeService(groceryRepo repositories.GroceryRepository, recipeRepo repositories.RecipeRepository, apiKey string) *RecipeService {
	return &RecipeService{
		groceryRepo: groceryRepo,
		recipeRepo:  recipeRepo,
		groqClient:  groq.NewClient(groq.WithAPIKey(apiKey)),
	}
}

// GenerateRecipes generates recipes based on the user's groceries and cuisine preference
func (s *RecipeService) GenerateRecipes(userID uint, cuisinePreference string) ([]models.Recipe, error) {
	// Get user's groceries
	groceries, err := s.groceryRepo.FindAll(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groceries: %w", err)
	}

	// Prepare prompt for Groq API
	prompt := s.buildPrompt(groceries, cuisinePreference)

	// Call Groq API
	response, err := s.groqClient.CreateChatCompletion(groq.CompletionCreateParams{
		Model: "llama3-70b-8192",
		Messages: []groq.Message{
			{
				Role:    "system",
				Content: "You are a professional chef that generates recipes based on available ingredients.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   1500,
	})
	if err != nil {
		return nil, fmt.Errorf("groq API error: %w", err)
	}

	// Log the raw response for debugging
	fmt.Printf("Raw AI Response: %s\n", response.Choices[0].Message.Content)

	// Parse the response into Recipe structs
	recipes, err := s.parseAIResponse(response.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Save recipes to the database
	for _, recipe := range recipes {
		recipe.UserID = userID // Associate the recipe with the user
		if err := s.recipeRepo.Save(&recipe); err != nil {
			return nil, fmt.Errorf("failed to save recipe to database: %w", err)
		}
	}

	return recipes, nil
}

// buildPrompt creates a prompt for the Groq API based on the user's groceries and cuisine preference
func (s *RecipeService) buildPrompt(groceries []models.GroceryItem, cuisine string) string {
	var ingredients []string
	for _, item := range groceries {
		ingredients = append(ingredients, fmt.Sprintf("%.1f %s %s", item.Quantity, item.Unit, item.Name))
	}

	prompt := fmt.Sprintf(`Generate 3 detailed recipes using these ingredients:
%s

Important instructions:
1. Cuisine preference: %s
2. Return ONLY a valid JSON array without any additional text or markdown formatting
3. Each recipe must have these exact fields:
   - title (string)
   - ingredients (string with each ingredient on a new line)
   - instructions (string with each step on a new line)
   - prep_time (number in minutes)
   - cook_time (number in minutes)
   - servings (number)
   - difficulty (string: Easy, Medium, or Hard)
   - cuisine (string)

Example format:
[
  {
    "title": "Recipe Name",
    "ingredients": "1 cup flour\n2 eggs",
    "instructions": "Step 1: Mix ingredients\nStep 2: Bake",
    "prep_time": 10,
    "cook_time": 20,
    "servings": 4,
    "difficulty": "Easy",
    "cuisine": "Italian"
  }
]`, strings.Join(ingredients, "\n"), cuisine)

	return prompt
}

// parseAIResponse parses the AI response into a slice of Recipe structs
func (s *RecipeService) parseAIResponse(response string) ([]models.Recipe, error) {
	// Clean the response by removing markdown code blocks if present
	cleanedResponse := strings.TrimSpace(response)
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```json")
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSuffix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	// Parse the JSON into Recipe structs
	var recipes []models.Recipe
	if err := json.Unmarshal([]byte(cleanedResponse), &recipes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AI response: %w\nResponse content: %s", err, cleanedResponse)
	}

	// Validate the parsed recipes
	if len(recipes) == 0 {
		return nil, errors.New("no recipes found in AI response")
	}

	for i, recipe := range recipes {
		if recipe.Title == "" || recipe.Instructions == "" {
			return nil, fmt.Errorf("invalid recipe format at index %d", i)
		}
	}

	return recipes, nil
}

// GetAllRecipes retrieves all recipes for a user
func (s *RecipeService) GetAllRecipes(userID uint) ([]models.Recipe, error) {
	recipes, err := s.recipeRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes: %w", err)
	}
	return recipes, nil
}

// GetRecipeByID retrieves a specific recipe by ID for a user
func (s *RecipeService) GetRecipeByID(userID uint, recipeID uint) (*models.Recipe, error) {
	recipe, err := s.recipeRepo.FindByID(userID, recipeID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, models.ErrRecipeNotFound
		}
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}
	return recipe, nil
}
