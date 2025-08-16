package models

import (
    "time"
    "gorm.io/gorm"
)

// Feedback model stores recipe feedback and ratings
type Feedback struct {
    gorm.Model
    RecipeID    uint      `json:"recipe_id" gorm:"not null"`
    UserID      uint      `json:"user_id" gorm:"not null"`
    Comment     string    `json:"comment" gorm:"type:text"`
    Rating      int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
    IsHelpful   bool      `json:"is_helpful" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at"`

    // Relationships
    Recipe      Recipe    `json:"recipe" gorm:"foreignKey:RecipeID"`
    User        User      `json:"user" gorm:"foreignKey:UserID"`
}

// IsValidRating checks if the rating is within valid range (1-5)
func (f *Feedback) IsValidRating() bool {
    return f.Rating >= 1 && f.Rating <= 5
}