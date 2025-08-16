package controllers

import (
    "net/http"
    "strconv"
    "shei-deli/models"
    "shei-deli/config"
    "github.com/gin-gonic/gin"
)

// GetRecipeFeedback fetches all feedback for a specific recipe
func GetRecipeFeedback(c *gin.Context) {
    recipeID := c.Param("recipeId")
    
    // Check if recipe exists
    var recipe models.Recipe
    if err := config.DB.First(&recipe, recipeID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }
    
    var feedbacks []models.Feedback
    if err := config.DB.Preload("User").Where("recipe_id = ?", recipeID).Order("created_at DESC").Find(&feedbacks).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving feedback"})
        return
    }
    
    // Calculate average rating
    var avgRating float64
    var totalRatings int64
    config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipeID).Select("AVG(rating)").Scan(&avgRating)
    config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipeID).Count(&totalRatings)
    
    c.JSON(http.StatusOK, gin.H{
        "feedbacks":      feedbacks,
        "average_rating": avgRating,
        "total_ratings":  totalRatings,
    })
}

// AddFeedback adds feedback to a recipe
func AddFeedback(c *gin.Context) {
    var feedback models.Feedback
    if err := c.ShouldBindJSON(&feedback); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback format"})
        return
    }
    
    // Validate rating
    if !feedback.IsValidRating() {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
        return
    }
    
    // Check if recipe exists
    var recipe models.Recipe
    if err := config.DB.First(&recipe, feedback.RecipeID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }
    
    // For now, use a default user ID (in a real app, this would come from authentication)
    if feedback.UserID == 0 {
        feedback.UserID = 1 // Default to admin user
    }
    
    // Check if user exists
    var user models.User
    if err := config.DB.First(&user, feedback.UserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    if err := config.DB.Create(&feedback).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving feedback"})
        return
    }
    
    // Load relationships for response
    config.DB.Preload("User").Preload("Recipe").First(&feedback, feedback.ID)
    
    c.JSON(http.StatusCreated, feedback)
}

// UpdateFeedback updates existing feedback
func UpdateFeedback(c *gin.Context) {
    id := c.Param("id")
    
    var feedback models.Feedback
    if err := config.DB.First(&feedback, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
        return
    }
    
    var updateData models.Feedback
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
        return
    }
    
    // Validate rating if provided
    if updateData.Rating != 0 && !updateData.IsValidRating() {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
        return
    }
    
    if err := config.DB.Model(&feedback).Updates(updateData).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating feedback"})
        return
    }
    
    // Load relationships for response
    config.DB.Preload("User").Preload("Recipe").First(&feedback, feedback.ID)
    
    c.JSON(http.StatusOK, feedback)
}

// DeleteFeedback deletes feedback
func DeleteFeedback(c *gin.Context) {
    id := c.Param("id")
    
    var feedback models.Feedback
    if err := config.DB.First(&feedback, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
        return
    }
    
    if err := config.DB.Delete(&feedback).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting feedback"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Feedback deleted successfully"})
}

// GetTopRatedRecipes fetches recipes with highest average ratings
func GetTopRatedRecipes(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    var recipes []models.Recipe
    
    // Get recipes with their average ratings
    if err := config.DB.Preload("User").
        Select("recipes.*, AVG(feedbacks.rating) as average_rating").
        Joins("LEFT JOIN feedbacks ON recipes.id = feedbacks.recipe_id").
        Group("recipes.id").
        Having("COUNT(feedbacks.id) > 0").
        Order("average_rating DESC").
        Limit(limit).
        Find(&recipes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving top rated recipes"})
        return
    }
    
    // Calculate average ratings for response
    for i := range recipes {
        var avgRating float64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        recipes[i].AverageRating = avgRating
    }
    
    c.JSON(http.StatusOK, gin.H{
        "recipes": recipes,
        "limit":   limit,
    })
}
