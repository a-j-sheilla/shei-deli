package controllers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "shei-deli/models"
    "shei-deli/config"
    "github.com/gin-gonic/gin"
)

// Fetch recipes from database
func GetRecipes(c *gin.Context) {
    var recipes []models.Recipe
    if err := config.DB.Find(&recipes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving recipes from the database"})
        return
    }
    c.JSON(http.StatusOK, recipes)
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

// Add a new recipe to the database
func AddRecipe(c *gin.Context) {
    var newRecipe models.Recipe
    if err := c.ShouldBindJSON(&newRecipe); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
        return
    }

    if err := config.DB.Create(&newRecipe).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving recipe to the database"})
        return
    }

    c.JSON(http.StatusCreated, newRecipe)
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
