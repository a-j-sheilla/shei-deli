package models

import "gorm.io/gorm"

type Recipe struct {
    gorm.Model  
    Title       string  
    Ingredients string  
    Instructions string  
    Category    string  
}
