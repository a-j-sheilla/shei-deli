package models

import "gorm.io/gorm"

// Recipe model to store recipe data
type Recipe struct {
    gorm.Model
    Title        string
    Ingredients  string
    Instructions string
    Category     string
    APIRecipeID  int `gorm:"default:null"` // This can store recipe ID from Spoonacular API
}
