package controllers

import (
    "database/sql"
    "net/http"
    "sort"
    "strconv"
    "time"
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
            Image:       "light-meals.jpeg",
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
            Image:       "soups.jpeg",
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
            Image:       "pastries.jpeg",
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
        var avgRating sql.NullFloat64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", featuredRecipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        if avgRating.Valid {
            featuredRecipes[i].AverageRating = avgRating.Float64
        } else {
            featuredRecipes[i].AverageRating = 0.0
        }
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
        var avgRating sql.NullFloat64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        if avgRating.Valid {
            recipes[i].AverageRating = avgRating.Float64
        } else {
            recipes[i].AverageRating = 0.0
        }
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
    var avgRating sql.NullFloat64
    config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipe.ID).Select("AVG(rating)").Scan(&avgRating)
    if avgRating.Valid {
        recipe.AverageRating = avgRating.Float64
    } else {
        recipe.AverageRating = 0.0
    }

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

// FeaturedHandler serves featured recipes page (highly-rated and popular recipes)
func FeaturedHandler(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit := 12
    offset := (page - 1) * limit

    // Get all recipes with their feedbacks to calculate ratings
    var allRecipes []models.Recipe
    if err := config.DB.Preload("User").Preload("Feedbacks").Find(&allRecipes).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Title": "Error",
            "Error": "Failed to load recipes.",
        })
        return
    }

    // Calculate average ratings and filter featured recipes
    var featuredRecipes []models.Recipe
    for i := range allRecipes {
        var avgRating sql.NullFloat64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", allRecipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        if avgRating.Valid {
            allRecipes[i].AverageRating = avgRating.Float64
        } else {
            allRecipes[i].AverageRating = 0.0
        }

        // Featured criteria: rating >= 4.0 OR has 2+ feedbacks OR is recent (created in last 30 days)
        feedbackCount := len(allRecipes[i].Feedbacks)
        isRecent := allRecipes[i].CreatedAt.After(time.Now().AddDate(0, 0, -30))

        if allRecipes[i].AverageRating >= 4.0 || feedbackCount >= 2 || isRecent {
            featuredRecipes = append(featuredRecipes, allRecipes[i])
        }
    }

    // Sort featured recipes by rating (descending), then by feedback count, then by creation date
    sort.Slice(featuredRecipes, func(i, j int) bool {
        if featuredRecipes[i].AverageRating != featuredRecipes[j].AverageRating {
            return featuredRecipes[i].AverageRating > featuredRecipes[j].AverageRating
        }
        if len(featuredRecipes[i].Feedbacks) != len(featuredRecipes[j].Feedbacks) {
            return len(featuredRecipes[i].Feedbacks) > len(featuredRecipes[j].Feedbacks)
        }
        return featuredRecipes[i].CreatedAt.After(featuredRecipes[j].CreatedAt)
    })

    // Apply pagination to featured recipes
    totalFeatured := len(featuredRecipes)
    totalPages := (totalFeatured + limit - 1) / limit

    start := offset
    end := offset + limit
    if start > totalFeatured {
        start = totalFeatured
    }
    if end > totalFeatured {
        end = totalFeatured
    }

    paginatedRecipes := featuredRecipes[start:end]

    c.HTML(http.StatusOK, "featured.html", gin.H{
        "Title":       "Featured Recipes",
        "Recipes":     paginatedRecipes,
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
