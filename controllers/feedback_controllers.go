package controllers

import (
    "github.com/gin-gonic/gin"
)

type Feedback struct {
    RecipeID uint   `json:"recipe_id"`
    Comment  string `json:"comment"`
    Rating   int    `json:"rating"`  // 1-5 stars rating
}

var feedbacks []Feedback  // also to be replaced with DB call

func SubmitFeedback(c *gin.Context) {
    var newFeedback Feedback
    if err := c.BindJSON(&newFeedback); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    feedbacks = append(feedbacks, newFeedback)
    c.JSON(200, gin.H{"message": "Feedback submitted successfully!"})
}
