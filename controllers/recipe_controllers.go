package controllers

import (
    "github.com/gin-gonic/gin"
    "shei-deli/models"
)

var recipes []models.Recipe  // to be replaced with actual DB calls

func ShowHomePage(c *gin.Context) {
    c.HTML(200, "index.html", gin.H{
        "title": "Shei-Deli Recipe Application",
    })
}

func ShowRecipesByCategory(c *gin.Context) {
    category := c.Param("category")
    filteredRecipes := []models.Recipe{}
    
    for _, recipe := range recipes {
        if recipe.Category == category {
            filteredRecipes = append(filteredRecipes, recipe)
        }
    }
    
    c.JSON(200, filteredRecipes)
}

func UploadRecipe(c *gin.Context) {
    var newRecipe models.Recipe
    if err := c.BindJSON(&newRecipe); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    recipes = append(recipes, newRecipe)
    c.JSON(200, gin.H{"message": "Recipe uploaded successfully!"})
}
