package models

import (
    "gorm.io/gorm"
)

// RecipeCategory defines the available recipe categories
type RecipeCategory string

const (
    PlantBasedMeals RecipeCategory = "plant_based_meals"
    KidsMeals       RecipeCategory = "kids_meals"
    LightMeals      RecipeCategory = "light_meals"
    HeartyMeals     RecipeCategory = "hearty_meals"
    MeatStews       RecipeCategory = "meat_stews"
    VeggieStews     RecipeCategory = "veggie_stews"
    SeafoodStews    RecipeCategory = "seafood_stews"
    FusionStews     RecipeCategory = "fusion_stews"
    Soups           RecipeCategory = "soups"
    Drinks          RecipeCategory = "drinks"
    Pastries        RecipeCategory = "pastries"
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
    case PlantBasedMeals:
        return "Plant-Based Meals"
    case KidsMeals:
        return "Kids' Meals"
    case LightMeals:
        return "Light Meals (Weight Loss)"
    case HeartyMeals:
        return "Hearty Meals (Weight Gain)"
    case MeatStews:
        return "Meat Stews"
    case VeggieStews:
        return "Veggie Stews"
    case SeafoodStews:
        return "Seafood & Fish Stews"
    case FusionStews:
        return "Fusion Stews"
    case Soups:
        return "Soups"
    case Drinks:
        return "Drinks"
    case Pastries:
        return "Pastries"
    default:
        return string(r)
    }
}

// IsValidCategory checks if the category is valid
func IsValidCategory(category string) bool {
    switch RecipeCategory(category) {
    case PlantBasedMeals, KidsMeals, LightMeals, HeartyMeals, MeatStews, VeggieStews, SeafoodStews, FusionStews, Soups, Drinks, Pastries:
        return true
    default:
        return false
    }
}
