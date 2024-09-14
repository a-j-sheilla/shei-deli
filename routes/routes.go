package routes

import (
    "github.com/gin-gonic/gin"
    "shei-deli/controllers"
)

func SetupRoutes(router *gin.Engine) {
    router.LoadHTMLGlob("templates/*")
    
    // Recipe routes
    router.GET("/", controllers.ShowHomePage)
    router.GET("/recipes/:category", controllers.ShowRecipesByCategory)
    router.POST("/recipes/upload", controllers.UploadRecipe)
    
    // Feedback routes
    router.POST("/recipes/:id/feedback", controllers.SubmitFeedback)
}
