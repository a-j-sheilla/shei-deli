package routes

import (
    "github.com/gin-gonic/gin"
    "shei-deli/controllers"
)

// SetupRoutes sets up API routes for the application
func SetupRoutes() *gin.Engine {
    router := gin.Default()

    // Serve static files
    router.Static("/static", "./static")
    router.Static("/images", "./images")

    // Web routes (HTML pages)
    router.GET("/", controllers.HomeHandler)
    router.GET("/category/:category", controllers.CategoryHandler)
    router.GET("/recipe/:id", controllers.RecipeHandler)
    router.GET("/add-recipe", controllers.AddRecipeHandler)
    router.GET("/register", controllers.RegisterHandler)
    router.GET("/featured", controllers.FeaturedHandler)
    router.GET("/about", controllers.AboutHandler)

    // API version 1 group
    v1 := router.Group("/api/v1")
    {
        // Recipe routes
        recipes := v1.Group("/recipes")
        {
            recipes.GET("", controllers.GetRecipes)                           // Get all recipes with optional category filter
            recipes.POST("", controllers.AddRecipe)                          // Add a new recipe
            recipes.GET("/:id", controllers.GetRecipeByID)                   // Get recipe by ID
            recipes.PUT("/:id", controllers.UpdateRecipe)                    // Update recipe
            recipes.DELETE("/:id", controllers.DeleteRecipe)                 // Delete recipe
            recipes.GET("/category/:category", controllers.GetRecipesByCategory) // Get recipes by category
            recipes.GET("/top-rated", controllers.GetTopRatedRecipes)        // Get top rated recipes
            recipes.GET("/search", controllers.GetSpoonacularRecipes)        // Search recipes using Spoonacular API
        }

        // Feedback routes
        feedback := v1.Group("/feedback")
        {
            feedback.POST("", controllers.AddFeedback)                       // Add feedback for a recipe
            feedback.GET("/recipe/:recipeId", controllers.GetRecipeFeedback) // Get all feedback for a recipe
            feedback.PUT("/:id", controllers.UpdateFeedback)                 // Update feedback
            feedback.DELETE("/:id", controllers.DeleteFeedback)              // Delete feedback
        }

        // User routes
        users := v1.Group("/users")
        {
            users.POST("/register", controllers.RegisterUser)               // Register new user
            users.POST("/login", controllers.LoginUser)                     // User login
            users.GET("", controllers.GetAllUsers)                          // Get all users (admin)
            users.GET("/:id", controllers.GetUserProfile)                   // Get user profile
            users.PUT("/:id", controllers.UpdateUserProfile)                // Update user profile
            users.GET("/:id/recipes", controllers.GetUserRecipes)           // Get user's recipes
        }



        // Category routes (for getting category information)
        categories := v1.Group("/categories")
        {
            categories.GET("", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "categories": []gin.H{
                        {"key": "plant_based_meals", "name": "Plant-Based Meals"},
                        {"key": "kids_meals", "name": "Kids' Meals"},
                        {"key": "light_meals", "name": "Light Meals (Weight Loss)"},
                        {"key": "hearty_meals", "name": "Hearty Meals (Weight Gain)"},
                        {"key": "meat_stews", "name": "Meat Stews"},
                        {"key": "veggie_stews", "name": "Veggie Stews"},
                        {"key": "seafood_stews", "name": "Seafood & Fish Stews"},
                        {"key": "fusion_stews", "name": "Fusion Stews"},
                        {"key": "soups", "name": "Soups"},
                        {"key": "drinks", "name": "Drinks"},
                        {"key": "pastries", "name": "Pastries"},
                    },
                })
            })
        }

        // External API integration routes
        external := v1.Group("/external")
        {
            external.GET("/recipes/:category", controllers.SearchExternalRecipes)
            external.GET("/api-mapping", controllers.GetAPIMappingInfo)
        }
    }

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "healthy",
            "service": "shei-deli-api",
        })
    })

    return router
}
