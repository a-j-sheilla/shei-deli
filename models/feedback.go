package models

import "gorm.io/gorm"


// stores the recipes feedback data
type Feedback struct {
	gorm.Model
	RecipeID uint
	Comment string
	Rating int
}