package models

import (
    "time"
    "gorm.io/gorm"
)

// User model stores user data for community features
type User struct {
    gorm.Model
    Username    string    `json:"username" gorm:"uniqueIndex;not null"`
    Email       string    `json:"email" gorm:"uniqueIndex;not null"`
    Password    string    `json:"-" gorm:"not null"` // Hidden from JSON responses
    FirstName   string    `json:"first_name"`
    LastName    string    `json:"last_name"`
    Bio         string    `json:"bio"`
    AvatarURL   string    `json:"avatar_url"`
    IsActive    bool      `json:"is_active" gorm:"default:true"`
    JoinedAt    time.Time `json:"joined_at" gorm:"autoCreateTime"`
    
    // Relationships
    Recipes     []Recipe   `json:"recipes" gorm:"foreignKey:UserID"`
    Feedbacks   []Feedback `json:"feedbacks" gorm:"foreignKey:UserID"`
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
    if u.FirstName != "" && u.LastName != "" {
        return u.FirstName + " " + u.LastName
    }
    return u.Username
}

// GetDisplayName returns the best available display name
func (u User) GetDisplayName() string {
    if u.FirstName != "" {
        return u.FirstName
    }
    return u.Username
}
