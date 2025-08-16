package models

import (
    "time"
    "gorm.io/gorm"
)

// RecipeCategory defines the available recipe categories
type RecipeCategory string

const (
    VeganMeals      RecipeCategory = "vegan_meals"
    KidsMeals       RecipeCategory = "kids_meals"
    WeightLossMeals RecipeCategory = "weight_loss_meals"
    WeightGainMeals RecipeCategory = "weight_gain_meals"
)

// Recipe model stores recipe data
type Recipe struct {
    gorm.Model
    Title           string         `json:"title" gorm:"not null"`
    Description     string         `json:"description"`
    Ingredients     string         `json:"ingredients" gorm:"type:text;not null"`
    Instructions    string         `json:"instructions" gorm:"type:text;not null"`
    Category        RecipeCategory `json:"category" gorm:"not null"`
    PrepTime        int            `json:"prep_time"` // in minutes
    CookTime        int            `json:"cook_time"` // in minutes
    Servings        int            `json:"servings"`
    Difficulty      string         `json:"difficulty"` // Easy, Medium, Hard
    ImageURL        string         `json:"image_url"`
    UserID          uint           `json:"user_id" gorm:"not null"` // Foreign key to User
    User            User           `json:"user" gorm:"foreignKey:UserID"`
    Feedbacks       []Feedback     `json:"feedbacks" gorm:"foreignKey:RecipeID"`
    AverageRating   float64        `json:"average_rating" gorm:"-"` // Calculated field
    APIRecipeID     *int           `json:"api_recipe_id" gorm:"default:null"` // stores the recipe ID from Spoonacular API
}

// GetCategoryDisplayName returns a human-readable category name
func (r RecipeCategory) GetDisplayName() string {
    switch r {
    case VeganMeals:
        return "Vegan Meals"
    case KidsMeals:
        return "Kids' Meals"
    case WeightLossMeals:
        return "Weight Loss Meals"
    case WeightGainMeals:
        return "Weight Gain Meals"
    default:
        return string(r)
    }
}

// IsValidCategory checks if the category is valid
func IsValidCategory(category string) bool {
    switch RecipeCategory(category) {
    case VeganMeals, KidsMeals, WeightLossMeals, WeightGainMeals:
        return true
    default:
        return false
    }
}
