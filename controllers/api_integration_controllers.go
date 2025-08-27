package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
    "shei-deli/models"
    "github.com/gin-gonic/gin"
)

// API Configuration
const (
    SpoonacularAPIKey = "095a1d67a1ac14b20872de17b6b95749" // Spoonacular API key
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

// SpoonacularRecipe represents a recipe from Spoonacular API
type SpoonacularRecipe struct {
    ID                int     `json:"id"`
    Title             string  `json:"title"`
    Summary           string  `json:"summary"`
    Image             string  `json:"image"`
    ReadyInMinutes    int     `json:"readyInMinutes"`
    Servings          int     `json:"servings"`
    SourceName        string  `json:"sourceName"`
    SourceUrl         string  `json:"sourceUrl"`
    SpoonacularScore  float64 `json:"spoonacularScore"`
    HealthScore       float64 `json:"healthScore"`
}

// SpoonacularSearchResponse represents the search response from Spoonacular
type SpoonacularSearchResponse struct {
    Results      []SpoonacularRecipe `json:"results"`
    Offset       int                 `json:"offset"`
    Number       int                 `json:"number"`
    TotalResults int                 `json:"totalResults"`
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

// fetchSpoonacularRecipes makes actual API calls to Spoonacular
func fetchSpoonacularRecipes(spoonacularConfig SpoonacularParams, limit int) ([]ExternalRecipe, error) {
    // Build API URL
    baseURL := SpoonacularBaseURL + "/complexSearch"
    params := url.Values{}

    // Set API key and basic parameters
    params.Set("apiKey", SpoonacularAPIKey)
    params.Set("number", strconv.Itoa(limit))
    params.Set("addRecipeInformation", "true")
    params.Set("fillIngredients", "false")
    params.Set("instructionsRequired", "true")

    // Add category-specific parameters
    if spoonacularConfig.Query != "" {
        params.Set("query", spoonacularConfig.Query)
    }
    if spoonacularConfig.Diet != "" {
        params.Set("diet", spoonacularConfig.Diet)
    }
    if spoonacularConfig.Type != "" {
        params.Set("type", spoonacularConfig.Type)
    }
    if spoonacularConfig.Cuisine != "" {
        params.Set("cuisine", spoonacularConfig.Cuisine)
    }
    if spoonacularConfig.MaxCalories > 0 {
        params.Set("maxCalories", strconv.Itoa(spoonacularConfig.MaxCalories))
    }
    if spoonacularConfig.MinCalories > 0 {
        params.Set("minCalories", strconv.Itoa(spoonacularConfig.MinCalories))
    }
    if spoonacularConfig.MaxFat > 0 {
        params.Set("maxFat", strconv.Itoa(spoonacularConfig.MaxFat))
    }
    if spoonacularConfig.MinProtein > 0 {
        params.Set("minProtein", strconv.Itoa(spoonacularConfig.MinProtein))
    }

    // Build full URL
    fullURL := baseURL + "?" + params.Encode()

    // Create HTTP client with timeout
    client := &http.Client{
        Timeout: 15 * time.Second,
    }

    // Make the API request
    resp, err := client.Get(fullURL)
    if err != nil {
        return nil, fmt.Errorf("failed to make API request: %v", err)
    }
    defer resp.Body.Close()

    // Check response status
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
    }

    // Parse the response
    var apiResponse SpoonacularSearchResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
        return nil, fmt.Errorf("failed to parse API response: %v", err)
    }

    // Convert Spoonacular recipes to our ExternalRecipe format
    var externalRecipes []ExternalRecipe
    for _, recipe := range apiResponse.Results {
        // Clean HTML tags from summary
        cleanSummary := cleanHTMLTags(recipe.Summary)
        if len(cleanSummary) > 200 {
            cleanSummary = cleanSummary[:200] + "..."
        }

        externalRecipe := ExternalRecipe{
            ID:          strconv.Itoa(recipe.ID),
            Title:       recipe.Title,
            Description: cleanSummary,
            Image:       recipe.Image,
            ReadyTime:   recipe.ReadyInMinutes,
            Servings:    recipe.Servings,
            Source:      "Spoonacular",
            SourceURL:   recipe.SourceUrl,
            Ingredients: []string{}, // Will be populated if needed
            Instructions: "Visit source for full instructions",
        }

        // Ensure we have a valid source URL
        if externalRecipe.SourceURL == "" {
            externalRecipe.SourceURL = fmt.Sprintf("https://spoonacular.com/recipes/%s-%d",
                strings.ReplaceAll(strings.ToLower(recipe.Title), " ", "-"), recipe.ID)
        }

        externalRecipes = append(externalRecipes, externalRecipe)
    }

    return externalRecipes, nil
}

// cleanHTMLTags removes HTML tags from text
func cleanHTMLTags(text string) string {
    result := text
    // Simple HTML tag removal
    for strings.Contains(result, "<") && strings.Contains(result, ">") {
        start := strings.Index(result, "<")
        end := strings.Index(result[start:], ">")
        if end != -1 {
            result = result[:start] + result[start+end+1:]
        } else {
            break
        }
    }
    return strings.TrimSpace(result)
}

// getEnhancedMockRecipes returns realistic mock data for each category
func getEnhancedMockRecipes(category models.RecipeCategory, limit int) []ExternalRecipe {
    mockData := map[models.RecipeCategory][]ExternalRecipe{
        models.PlantBasedMeals: {
            {ID: "mock_pb_1", Title: "Quinoa Buddha Bowl", Description: "Nutritious bowl with quinoa, roasted vegetables, and tahini dressing", Image: "/images/plant-based.jpg", ReadyTime: 25, Servings: 2, Source: "Plant-Based Kitchen", SourceURL: "https://example.com/quinoa-bowl", Ingredients: []string{"Quinoa", "Sweet potato", "Chickpeas", "Tahini"}, Instructions: "Cook quinoa, roast vegetables, assemble bowl with tahini dressing"},
            {ID: "mock_pb_2", Title: "Lentil Curry", Description: "Spicy red lentil curry with coconut milk and aromatic spices", Image: "/images/plant-based.jpg", ReadyTime: 30, Servings: 4, Source: "Vegan Delights", SourceURL: "https://example.com/lentil-curry", Ingredients: []string{"Red lentils", "Coconut milk", "Curry spices", "Tomatoes"}, Instructions: "Simmer lentils with spices and coconut milk until tender"},
            {ID: "mock_pb_3", Title: "Avocado Toast Supreme", Description: "Gourmet avocado toast with hemp seeds and microgreens", Image: "/images/plant-based.jpg", ReadyTime: 10, Servings: 1, Source: "Healthy Eats", SourceURL: "https://example.com/avocado-toast", Ingredients: []string{"Sourdough bread", "Avocado", "Hemp seeds", "Microgreens"}, Instructions: "Toast bread, mash avocado, top with seeds and greens"},
        },
        models.KidsMeals: {
            {ID: "mock_km_1", Title: "Mini Cheese Quesadillas", Description: "Kid-friendly quesadillas with mild cheese and hidden vegetables", Image: "/images/kids-meals.jpg", ReadyTime: 15, Servings: 2, Source: "Family Kitchen", SourceURL: "https://example.com/mini-quesadillas", Ingredients: []string{"Tortillas", "Mild cheese", "Hidden veggies"}, Instructions: "Fill tortillas with cheese and veggies, cook until golden"},
            {ID: "mock_km_2", Title: "Chicken Nuggets", Description: "Homemade baked chicken nuggets that kids absolutely love", Image: "/images/kids-meals.jpg", ReadyTime: 25, Servings: 4, Source: "Kid-Approved", SourceURL: "https://example.com/chicken-nuggets", Ingredients: []string{"Chicken breast", "Breadcrumbs", "Eggs"}, Instructions: "Coat chicken in breadcrumbs and bake until crispy"},
            {ID: "mock_km_3", Title: "Mac and Cheese Cups", Description: "Individual mac and cheese portions baked in muffin cups", Image: "/images/kids-meals.jpg", ReadyTime: 20, Servings: 6, Source: "Fun Foods", SourceURL: "https://example.com/mac-cheese-cups", Ingredients: []string{"Pasta", "Cheese sauce", "Breadcrumbs"}, Instructions: "Mix pasta with cheese, bake in muffin tins"},
        },
        models.LightMeals: {
            {ID: "mock_lm_1", Title: "Greek Salad", Description: "Fresh Mediterranean salad with feta cheese and olives", Image: "/images/light-meals.jpg", ReadyTime: 10, Servings: 2, Source: "Healthy Living", SourceURL: "https://example.com/greek-salad", Ingredients: []string{"Cucumber", "Tomatoes", "Feta", "Olives"}, Instructions: "Chop vegetables, add feta and olives, dress with olive oil"},
            {ID: "mock_lm_2", Title: "Zucchini Noodles", Description: "Low-carb zucchini noodles with fresh pesto sauce", Image: "/images/light-meals.jpg", ReadyTime: 15, Servings: 2, Source: "Diet Kitchen", SourceURL: "https://example.com/zucchini-noodles", Ingredients: []string{"Zucchini", "Basil pesto", "Cherry tomatoes"}, Instructions: "Spiralize zucchini, toss with pesto and tomatoes"},
            {ID: "mock_lm_3", Title: "Cauliflower Rice Bowl", Description: "Nutritious cauliflower rice with grilled vegetables", Image: "/images/light-meals.jpg", ReadyTime: 20, Servings: 2, Source: "Clean Eating", SourceURL: "https://example.com/cauliflower-rice", Ingredients: []string{"Cauliflower", "Bell peppers", "Broccoli"}, Instructions: "Rice cauliflower, grill vegetables, combine in bowl"},
        },
        models.Drinks: {
            {ID: "mock_dr_1", Title: "Green Smoothie", Description: "Healthy green smoothie with spinach, banana, and mango", Image: "/images/drinks.jpg", ReadyTime: 5, Servings: 1, Source: "Smoothie Bar", SourceURL: "https://example.com/green-smoothie", Ingredients: []string{"Spinach", "Banana", "Mango", "Coconut water"}, Instructions: "Blend all ingredients until smooth"},
            {ID: "mock_dr_2", Title: "Iced Coffee", Description: "Refreshing iced coffee with vanilla and cream", Image: "/images/drinks.jpg", ReadyTime: 10, Servings: 1, Source: "Coffee Shop", SourceURL: "https://example.com/iced-coffee", Ingredients: []string{"Coffee", "Ice", "Vanilla", "Cream"}, Instructions: "Brew coffee, add ice and flavorings"},
            {ID: "mock_dr_3", Title: "Fruit Infused Water", Description: "Refreshing water infused with fresh fruits and herbs", Image: "/images/drinks.jpg", ReadyTime: 5, Servings: 4, Source: "Hydration Station", SourceURL: "https://example.com/infused-water", Ingredients: []string{"Water", "Cucumber", "Mint", "Lemon"}, Instructions: "Add fruits and herbs to water, let infuse"},
        },
        models.Pastries: {
            {ID: "mock_pa_1", Title: "Chocolate Croissants", Description: "Buttery croissants filled with rich dark chocolate", Image: "/images/pastries.jpg", ReadyTime: 180, Servings: 8, Source: "French Bakery", SourceURL: "https://example.com/chocolate-croissants", Ingredients: []string{"Puff pastry", "Dark chocolate", "Butter", "Egg wash"}, Instructions: "Wrap chocolate in pastry, proof, and bake until golden"},
            {ID: "mock_pa_2", Title: "Apple Turnovers", Description: "Flaky pastry turnovers filled with spiced apples", Image: "/images/pastries.jpg", ReadyTime: 45, Servings: 6, Source: "Pastry Chef", SourceURL: "https://example.com/apple-turnovers", Ingredients: []string{"Puff pastry", "Apples", "Cinnamon", "Sugar"}, Instructions: "Fill pastry with spiced apples, seal and bake"},
            {ID: "mock_pa_3", Title: "Blueberry Muffins", Description: "Fluffy muffins bursting with fresh blueberries", Image: "/images/pastries.jpg", ReadyTime: 30, Servings: 12, Source: "Bakehouse", SourceURL: "https://example.com/blueberry-muffins", Ingredients: []string{"Flour", "Blueberries", "Sugar", "Eggs"}, Instructions: "Mix batter, fold in blueberries, bake in muffin tins"},
        },
    }

    // Add default recipes for categories not explicitly defined
    defaultRecipes := []ExternalRecipe{
        {ID: "mock_def_1", Title: fmt.Sprintf("Delicious %s Recipe", category.GetDisplayName()), Description: fmt.Sprintf("A wonderful %s recipe from our curated collection", strings.ToLower(category.GetDisplayName())), Image: fmt.Sprintf("/images/%s.jpg", strings.ReplaceAll(strings.ToLower(category.GetDisplayName()), " ", "-")), ReadyTime: 30, Servings: 4, Source: "Recipe Collection", SourceURL: "https://example.com/recipe", Ingredients: []string{"Fresh ingredients", "Quality seasonings", "Love and care"}, Instructions: "Follow traditional cooking methods for best results"},
        {ID: "mock_def_2", Title: fmt.Sprintf("Classic %s Dish", category.GetDisplayName()), Description: fmt.Sprintf("Traditional %s preparation with modern touches", strings.ToLower(category.GetDisplayName())), Image: fmt.Sprintf("/images/%s.jpg", strings.ReplaceAll(strings.ToLower(category.GetDisplayName()), " ", "-")), ReadyTime: 45, Servings: 6, Source: "Traditional Kitchen", SourceURL: "https://example.com/classic", Ingredients: []string{"Traditional ingredients", "Authentic spices", "Time-tested methods"}, Instructions: "Prepare using time-honored techniques"},
    }

    recipes, exists := mockData[category]
    if !exists {
        recipes = defaultRecipes
    }

    // Limit results
    if limit > 0 && limit < len(recipes) {
        recipes = recipes[:limit]
    }

    return recipes
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

    // Get limit from query parameter (default: 12, max: 50)
    limit := 12
    if l := c.Query("limit"); l != "" {
        if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
            limit = parsed
        }
    }

    // Fetch recipes from Spoonacular API
    externalRecipes, err := fetchSpoonacularRecipes(mapping.Spoonacular, limit)
    if err != nil {
        // If API call fails, return enhanced mock data as fallback
        mockRecipes := getEnhancedMockRecipes(category, limit)

        c.JSON(http.StatusOK, gin.H{
            "category":         category.GetDisplayName(),
            "api_mapping":      mapping,
            "external_recipes": mockRecipes,
            "note":            "Using enhanced mock data. Real API integration ready - API key may need activation.",
            "source":          "enhanced_mock",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "category":         category.GetDisplayName(),
        "api_mapping":      mapping,
        "external_recipes": externalRecipes,
        "count":           len(externalRecipes),
        "source":          "spoonacular_api",
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
