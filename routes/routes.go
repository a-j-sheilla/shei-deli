package routes

import (
    "github.com/gin-gonic/gin"
    "shei-deli/controllers"
)

// SetupRoutes sets up API routes for the application
func SetupRoutes() *gin.Engine {
    router := gin.Default()

    // Recipe routes
    router.GET("/recipes", controllers.GetRecipes)                     // Get all recipes from the database
    router.GET("/recipes/search", controllers.GetSpoonacularRecipes)    // Search recipes using Spoonacular API
    router.POST("/recipes", controllers.AddRecipe)                     // Add a new recipe

    // Feedback routes
    router.POST("/feedback", controllers.AddFeedback)                  // Add feedback for a recipe

    return router
}
