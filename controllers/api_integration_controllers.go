package controllers

import (
    "fmt"
    "net/http"
    "shei-deli/models"
    "github.com/gin-gonic/gin"
)

// API Configuration
const (
    SpoonacularAPIKey = " bd6a8785a56f4fcbb5438f332dd328a6 " // Replace with actual API key
    EdamamAppID       = "YOUR_EDAMAM_APP_ID"       // Replace with actual App ID
    EdamamAppKey      = "YOUR_EDAMAM_APP_KEY"      // Replace with actual App Key
)

// API URLs
const (
    SpoonacularBaseURL = "https://api.spoonacular.com/recipes"
    EdamamBaseURL      = "https://api.edamam.com/search"
    TheMealDBBaseURL   = "https://www.themealdb.com/api/json/v1/1"
)

// ExternalRecipe represents a recipe from external APIs
type ExternalRecipe struct {
    ID          string  `json:"id"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Image       string  `json:"image"`
    ReadyTime   int     `json:"ready_time"`
    Servings    int     `json:"servings"`
    Source      string  `json:"source"`
    SourceURL   string  `json:"source_url"`
    Ingredients []string `json:"ingredients"`
    Instructions string `json:"instructions"`
}

// CategoryAPIMapping defines how each category maps to external APIs
type CategoryAPIMapping struct {
    Category    models.RecipeCategory
    Spoonacular SpoonacularParams
    Edamam      EdamamParams
    TheMealDB   TheMealDBParams
}

type SpoonacularParams struct {
    Query       string
    Diet        string
    MaxCalories int
    MinCalories int
    MaxFat      int
    MinProtein  int
    Cuisine     string
    Type        string
}

type EdamamParams struct {
    Query       string
    Diet        string
    Health      string
    CuisineType string
    MealType    string
    CaloriesMin int
    CaloriesMax int
}

type TheMealDBParams struct {
    Category string
    Area     string
    Query    string
}

// GetCategoryAPIMapping returns API parameters for each category
func GetCategoryAPIMapping() map[models.RecipeCategory]CategoryAPIMapping {
    return map[models.RecipeCategory]CategoryAPIMapping{
        models.PlantBasedMeals: {
            Category: models.PlantBasedMeals,
            Spoonacular: SpoonacularParams{
                Diet: "vegan,vegetarian",
                Type: "main course",
            },
            Edamam: EdamamParams{
                Diet:   "vegan,vegetarian",
                Health: "vegan,vegetarian",
            },
        },
        models.KidsMeals: {
            Category: models.KidsMeals,
            Spoonacular: SpoonacularParams{
                Query:       "kid friendly",
                MaxCalories: 500,
                Type:        "main course",
            },
            Edamam: EdamamParams{
                Query:       "kid friendly",
                CaloriesMax: 500,
                Health:      "low-sugar",
            },
        },
        models.LightMeals: {
            Category: models.LightMeals,
            Spoonacular: SpoonacularParams{
                MaxCalories: 400,
                MaxFat:      15,
                Type:        "main course",
            },
            Edamam: EdamamParams{
                CaloriesMax: 400,
                Diet:        "low-fat",
                Health:      "low-fat",
            },
        },
        models.HeartyMeals: {
            Category: models.HeartyMeals,
            Spoonacular: SpoonacularParams{
                MinCalories: 600,
                MinProtein:  25,
                Type:        "main course",
            },
            Edamam: EdamamParams{
                CaloriesMin: 600,
                Health:      "high-protein",
            },
        },
        models.MeatStews: {
            Category: models.MeatStews,
            Spoonacular: SpoonacularParams{
                Query: "beef stew,chicken stew,lamb stew",
                Type:  "main course",
            },
            TheMealDB: TheMealDBParams{
                Query: "beef,chicken,lamb",
            },
        },
        models.VeggieStews: {
            Category: models.VeggieStews,
            Spoonacular: SpoonacularParams{
                Query: "vegetable stew",
                Diet:  "vegetarian",
                Type:  "main course",
            },
            Edamam: EdamamParams{
                Query: "vegetable stew",
                Diet:  "vegetarian",
            },
        },
        models.SeafoodStews: {
            Category: models.SeafoodStews,
            Spoonacular: SpoonacularParams{
                Query: "fish stew,seafood stew",
                Type:  "main course",
            },
            TheMealDB: TheMealDBParams{
                Category: "Seafood",
            },
        },
        models.FusionStews: {
            Category: models.FusionStews,
            Spoonacular: SpoonacularParams{
                Query:   "stew,curry",
                Cuisine: "African,Asian,European,Latin American",
                Type:    "main course",
            },
            Edamam: EdamamParams{
                Query:       "stew,curry",
                CuisineType: "Asian,European,American",
            },
        },
        models.Soups: {
            Category: models.Soups,
            Spoonacular: SpoonacularParams{
                Query: "soup",
                Type:  "soup",
            },
            TheMealDB: TheMealDBParams{
                Query: "soup",
            },
        },
        models.Drinks: {
            Category: models.Drinks,
            Spoonacular: SpoonacularParams{
                Query: "smoothie,juice,tea",
                Type:  "drink",
            },
            Edamam: EdamamParams{
                Query:    "smoothie,juice",
                MealType: "snack",
            },
        },
        models.Pastries: {
            Category: models.Pastries,
            Spoonacular: SpoonacularParams{
                Query: "pastry,cake,cookie,pie,bread",
                Type:  "dessert,bread",
            },
            TheMealDB: TheMealDBParams{
                Category: "Dessert",
            },
        },
    }
}

// SearchExternalRecipes searches external APIs for recipes based on category
func SearchExternalRecipes(c *gin.Context) {
    categoryStr := c.Param("category")
    
    if !models.IsValidCategory(categoryStr) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
        return
    }
    
    category := models.RecipeCategory(categoryStr)
    mapping := GetCategoryAPIMapping()[category]
    
    // For demo purposes, return mock data with API mapping information
    // In production, you would make actual API calls here
    
    mockRecipes := []ExternalRecipe{
        {
            ID:          "ext_1",
            Title:       fmt.Sprintf("External %s Recipe", mapping.Category.GetDisplayName()),
            Description: fmt.Sprintf("A delicious recipe from external APIs for %s", mapping.Category.GetDisplayName()),
            Image:       "https://via.placeholder.com/300x200",
            ReadyTime:   30,
            Servings:    4,
            Source:      "Spoonacular",
            SourceURL:   "https://spoonacular.com",
            Ingredients: []string{"Ingredient 1", "Ingredient 2", "Ingredient 3"},
            Instructions: "1. Prepare ingredients\n2. Cook according to recipe\n3. Serve and enjoy",
        },
    }
    
    c.JSON(http.StatusOK, gin.H{
        "category":        category.GetDisplayName(),
        "api_mapping":     mapping,
        "external_recipes": mockRecipes,
        "note":           "This is a demo response. In production, this would fetch from actual APIs.",
    })
}

// GetAPIMappingInfo returns the API mapping configuration for all categories
func GetAPIMappingInfo(c *gin.Context) {
    mappings := GetCategoryAPIMapping()
    
    response := make(map[string]interface{})
    for category, mapping := range mappings {
        response[string(category)] = gin.H{
            "category_name": category.GetDisplayName(),
            "spoonacular":   mapping.Spoonacular,
            "edamam":        mapping.Edamam,
            "themealdb":     mapping.TheMealDB,
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "api_mappings": response,
        "description":  "API mapping configuration for Shei-deli recipe categories",
    })
}
