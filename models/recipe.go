package models

import "gorm.io/gorm"

// this mdel stores recipe data
type Recipe struct {
    gorm.Model
    Title        string
    Ingredients  string
    Instructions string
    Category     string
    APIRecipeID  int `gorm:"default:null"` // stores the recipe ID from Spoonacular API
}
