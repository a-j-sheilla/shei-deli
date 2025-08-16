package controllers

import (
    "net/http"
    "strconv"
    "shei-deli/models"
    "shei-deli/config"
    "github.com/gin-gonic/gin"
)

// Category represents a recipe category for the UI
type Category struct {
    Key         string
    Name        string
    Description string
    Image       string
}

// Stats represents application statistics
type Stats struct {
    TotalRecipes  int64
    TotalUsers    int64
    TotalFeedback int64
}

// GetCategories returns all available categories
func getCategories() []Category {
    return []Category{
        {
            Key:         "plant_based_meals",
            Name:        "Plant-Based Meals",
            Description: "Vegan/vegetarian options (no animal products)",
            Image:       "vegan.jpeg",
        },
        {
            Key:         "kids_meals",
            Name:        "Kids' Meals",
            Description: "Fun, simple, and nutritious meals for children",
            Image:       "kids-meals.jpeg",
        },
        {
            Key:         "light_meals",
            Name:        "Light Meals (Weight Loss)",
            Description: "Low-calorie, balanced recipes",
            Image:       "hearty-meals.jpeg",
        },
        {
            Key:         "hearty_meals",
            Name:        "Hearty Meals (Weight Gain)",
            Description: "High-calorie, energy-packed recipes",
            Image:       "hearty-meals.jpeg",
        },
        {
            Key:         "meat_stews",
            Name:        "Meat Stews",
            Description: "Beef, chicken, goat, lamb, and other meat-based stews",
            Image:       "stews.jpeg",
        },
        {
            Key:         "veggie_stews",
            Name:        "Veggie Stews",
            Description: "Lentil, bean, mushroom, and vegetable stews",
            Image:       "vegetable-stews.jpeg",
        },
        {
            Key:         "seafood_stews",
            Name:        "Seafood & Fish Stews",
            Description: "Fish stews, seafood mixes, and ocean-inspired flavors",
            Image:       "fish&sea-food.jpeg",
        },
        {
            Key:         "fusion_stews",
            Name:        "Fusion Stews",
            Description: "Cultural and traditional varieties (e.g., goulash, curries)",
            Image:       "fusion.jpeg",
        },
        {
            Key:         "soups",
            Name:        "Soups",
            Description: "Warm, comforting soups",
            Image:       "stews.jpeg",
        },
        {
            Key:         "drinks",
            Name:        "Drinks",
            Description: "Smoothies, juices, teas, and other beverages",
            Image:       "drinks&smoothies.jpeg",
        },
        {
            Key:         "pastries",
            Name:        "Pastries",
            Description: "Baked goods such as cakes, cookies, pies, and breads",
            Image:       "fusion.jpeg", // Using fusion image as placeholder
        },
    }
}

// HomeHandler serves the home page
func HomeHandler(c *gin.Context) {
    // Get featured recipes (top rated)
    var featuredRecipes []models.Recipe
    config.DB.Preload("User").Preload("Feedbacks").
        Joins("LEFT JOIN feedbacks ON recipes.id = feedbacks.recipe_id").
        Group("recipes.id").
        Order("AVG(feedbacks.rating) DESC").
        Limit(6).
        Find(&featuredRecipes)

    // Calculate average ratings
    for i := range featuredRecipes {
        var avgRating float64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", featuredRecipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        featuredRecipes[i].AverageRating = avgRating
    }

    // Get stats
    var stats Stats
    config.DB.Model(&models.Recipe{}).Count(&stats.TotalRecipes)
    config.DB.Model(&models.User{}).Count(&stats.TotalUsers)
    config.DB.Model(&models.Feedback{}).Count(&stats.TotalFeedback)

    c.HTML(http.StatusOK, "index.html", gin.H{
        "Title":           "Home",
        "Categories":      getCategories(),
        "FeaturedRecipes": featuredRecipes,
        "Stats":           stats,
    })
}

// CategoryHandler serves the category page
func CategoryHandler(c *gin.Context) {
    categoryKey := c.Param("category")
    
    // Find category info
    var categoryInfo Category
    categories := getCategories()
    for _, cat := range categories {
        if cat.Key == categoryKey {
            categoryInfo = cat
            break
        }
    }
    
    if categoryInfo.Key == "" {
        c.HTML(http.StatusNotFound, "error.html", gin.H{
            "Title": "Category Not Found",
            "Error": "The requested category was not found.",
        })
        return
    }

    // Get recipes for this category
    var recipes []models.Recipe
    if err := config.DB.Preload("User").Preload("Feedbacks").Where("category = ?", categoryKey).Find(&recipes).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Title": "Error",
            "Error": "Failed to load recipes for this category.",
        })
        return
    }

    // Calculate average ratings
    for i := range recipes {
        var avgRating float64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        recipes[i].AverageRating = avgRating
    }

    c.HTML(http.StatusOK, "category.html", gin.H{
        "Title":               categoryInfo.Name,
        "CategoryName":        categoryInfo.Name,
        "CategoryDescription": categoryInfo.Description,
        "CategoryKey":         categoryKey,
        "Recipes":             recipes,
    })
}

// RecipeHandler serves the recipe detail page
func RecipeHandler(c *gin.Context) {
    id := c.Param("id")
    
    var recipe models.Recipe
    if err := config.DB.Preload("User").Preload("Feedbacks.User").First(&recipe, id).Error; err != nil {
        c.HTML(http.StatusNotFound, "error.html", gin.H{
            "Title": "Recipe Not Found",
            "Error": "The requested recipe was not found.",
        })
        return
    }

    // Calculate average rating
    var avgRating float64
    config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipe.ID).Select("AVG(rating)").Scan(&avgRating)
    recipe.AverageRating = avgRating

    c.HTML(http.StatusOK, "recipe.html", gin.H{
        "Title":  recipe.Title,
        "Recipe": recipe,
    })
}

// AddRecipeHandler serves the add recipe form
func AddRecipeHandler(c *gin.Context) {
    selectedCategory := c.Query("category")
    
    c.HTML(http.StatusOK, "add-recipe.html", gin.H{
        "Title":            "Add Recipe",
        "Categories":       getCategories(),
        "SelectedCategory": selectedCategory,
    })
}

// RegisterHandler serves the user registration form
func RegisterHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "register.html", gin.H{
        "Title": "Join Community",
    })
}

// AllRecipesHandler serves all recipes page
func AllRecipesHandler(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit := 12
    offset := (page - 1) * limit

    var recipes []models.Recipe
    if err := config.DB.Preload("User").Preload("Feedbacks").
        Offset(offset).Limit(limit).Find(&recipes).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Title": "Error",
            "Error": "Failed to load recipes.",
        })
        return
    }

    // Calculate average ratings
    for i := range recipes {
        var avgRating float64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        recipes[i].AverageRating = avgRating
    }

    var totalRecipes int64
    config.DB.Model(&models.Recipe{}).Count(&totalRecipes)
    totalPages := (int(totalRecipes) + limit - 1) / limit

    c.HTML(http.StatusOK, "all-recipes.html", gin.H{
        "Title":       "All Recipes",
        "Recipes":     recipes,
        "CurrentPage": page,
        "TotalPages":  totalPages,
        "HasNext":     page < totalPages,
        "HasPrev":     page > 1,
    })
}

// AboutHandler serves the about page
func AboutHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "about.html", gin.H{
        "Title": "About Shei-deli",
    })
}
