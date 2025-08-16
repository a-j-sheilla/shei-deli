package controllers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "shei-deli/models"
    "shei-deli/config"
    "github.com/gin-gonic/gin"
)

// GetRecipes fetches all recipes from database with optional category filtering
func GetRecipes(c *gin.Context) {
    var recipes []models.Recipe
    query := config.DB.Preload("User").Preload("Feedbacks")

    // Filter by category if provided
    category := c.Query("category")
    if category != "" {
        if !models.IsValidCategory(category) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
            return
        }
        query = query.Where("category = ?", category)
    }

    // Add pagination
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset := (page - 1) * limit

    if err := query.Offset(offset).Limit(limit).Find(&recipes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving recipes from the database"})
        return
    }

    // Calculate average ratings for each recipe
    for i := range recipes {
        var avgRating float64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        recipes[i].AverageRating = avgRating
    }

    c.JSON(http.StatusOK, gin.H{
        "recipes": recipes,
        "page":    page,
        "limit":   limit,
    })
}

// GetRecipesByCategory fetches recipes by specific category
func GetRecipesByCategory(c *gin.Context) {
    category := c.Param("category")

    if !models.IsValidCategory(category) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
        return
    }

    var recipes []models.Recipe
    if err := config.DB.Preload("User").Preload("Feedbacks").Where("category = ?", category).Find(&recipes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving recipes"})
        return
    }

    // Calculate average ratings
    for i := range recipes {
        var avgRating float64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        recipes[i].AverageRating = avgRating
    }

    c.JSON(http.StatusOK, gin.H{
        "category": models.RecipeCategory(category).GetDisplayName(),
        "recipes":  recipes,
    })
}

// GetRecipeByID fetches a single recipe by ID
func GetRecipeByID(c *gin.Context) {
    id := c.Param("id")

    var recipe models.Recipe
    if err := config.DB.Preload("User").Preload("Feedbacks.User").First(&recipe, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }

    // Calculate average rating
    var avgRating float64
    config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipe.ID).Select("AVG(rating)").Scan(&avgRating)
    recipe.AverageRating = avgRating

    c.JSON(http.StatusOK, recipe)
}

// Fetch recipes from Spoonacular API
func GetSpoonacularRecipes(c *gin.Context) {
    apiKey := "your_spoonacular_api_key" // replace with your Spoonacular API key
    query := c.Query("query") // Fetch search query from URL parameter

    spoonacularURL := fmt.Sprintf("https://api.spoonacular.com/recipes/complexSearch?query=%s&apiKey=%s", query, apiKey)

    resp, err := http.Get(spoonacularURL)
    if err != nil || resp.StatusCode != http.StatusOK {
        log.Println("Error fetching data from Spoonacular API:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recipes from Spoonacular"})
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var result map[string]interface{}
    json.Unmarshal(body, &result)

    c.JSON(http.StatusOK, result)
}

// AddRecipe creates a new recipe in the database
func AddRecipe(c *gin.Context) {
    var newRecipe models.Recipe
    if err := c.ShouldBindJSON(&newRecipe); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
        return
    }

    // Validate category
    if !models.IsValidCategory(string(newRecipe.Category)) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
        return
    }

    // Validate required fields
    if newRecipe.Title == "" || newRecipe.Ingredients == "" || newRecipe.Instructions == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Title, ingredients, and instructions are required"})
        return
    }

    // For now, use a default user ID (in a real app, this would come from authentication)
    if newRecipe.UserID == 0 {
        newRecipe.UserID = 1 // Default to admin user
    }

    if err := config.DB.Create(&newRecipe).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving recipe to the database"})
        return
    }

    // Load the user relationship for the response
    config.DB.Preload("User").First(&newRecipe, newRecipe.ID)

    c.JSON(http.StatusCreated, newRecipe)
}

// UpdateRecipe updates an existing recipe
func UpdateRecipe(c *gin.Context) {
    id := c.Param("id")

    var recipe models.Recipe
    if err := config.DB.First(&recipe, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }

    var updateData models.Recipe
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
        return
    }

    // Validate category if provided
    if updateData.Category != "" && !models.IsValidCategory(string(updateData.Category)) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
        return
    }

    if err := config.DB.Model(&recipe).Updates(updateData).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating recipe"})
        return
    }

    // Load relationships for response
    config.DB.Preload("User").First(&recipe, recipe.ID)

    c.JSON(http.StatusOK, recipe)
}

// DeleteRecipe deletes a recipe from the database
func DeleteRecipe(c *gin.Context) {
    id := c.Param("id")

    var recipe models.Recipe
    if err := config.DB.First(&recipe, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }

    if err := config.DB.Delete(&recipe).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting recipe"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

// Add feedback to a recipe
func AddFeedback(c *gin.Context) {
    var feedback models.Feedback
    if err := c.ShouldBindJSON(&feedback); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback format"})
        return
    }

    // Check if recipe exists
    var recipe models.Recipe
    if err := config.DB.First(&recipe, feedback.RecipeID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
        return
    }

    if err := config.DB.Create(&feedback).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving feedback"})
        return
    }

    c.JSON(http.StatusCreated, feedback)
}
