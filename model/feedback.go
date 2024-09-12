package models

import "gorm.io/gorm"

// Feedback model
type Feedback struct {
    gorm.Model
    RecipeID uint   
    Comment  string 
    Rating   int    
}
