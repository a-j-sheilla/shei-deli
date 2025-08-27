package controllers

import (
    "net/http"
    "shei-deli/models"
    "shei-deli/config"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

// UserRegistration represents the registration request
type UserRegistration struct {
    Username  string `json:"username" binding:"required"`
    Email     string `json:"email" binding:"required"`
    Password  string `json:"password" binding:"required"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Bio       string `json:"bio"`
}

// RegisterUser creates a new user account
func RegisterUser(c *gin.Context) {
    var regData UserRegistration
    if err := c.ShouldBindJSON(&regData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
        return
    }

    // Validate required fields
    if regData.Username == "" || regData.Email == "" || regData.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
        return
    }
    
    // Check if username or email already exists
    var existingUser models.User
    if err := config.DB.Where("username = ? OR email = ?", regData.Username, regData.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
        return
    }

    // Create user model
    user := models.User{
        Username:  regData.Username,
        Email:     regData.Email,
        Password:  string(hashedPassword),
        FirstName: regData.FirstName,
        LastName:  regData.LastName,
        Bio:       regData.Bio,
        IsActive:  true,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }
    
    // Remove password from response
    user.Password = ""
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "user":    user,
    })
}

// LoginUser authenticates a user (basic implementation)
func LoginUser(c *gin.Context) {
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
        return
    }
    
    // Find user by username or email
    var user models.User
    if err := config.DB.Where("username = ? OR email = ?", loginData.Username, loginData.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Check if user is active
    if !user.IsActive {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
        return
    }
    
    // Remove password from response
    user.Password = ""
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user":    user,
        // In a real application, you would return a JWT token here
    })
}

// GetUserProfile fetches user profile information
func GetUserProfile(c *gin.Context) {
    id := c.Param("id")
    
    var user models.User
    if err := config.DB.Preload("Recipes").Preload("Feedbacks").First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    // Remove password from response
    user.Password = ""
    
    // Calculate user statistics
    var recipeCount int64
    var feedbackCount int64
    config.DB.Model(&models.Recipe{}).Where("user_id = ?", user.ID).Count(&recipeCount)
    config.DB.Model(&models.Feedback{}).Where("user_id = ?", user.ID).Count(&feedbackCount)
    
    c.JSON(http.StatusOK, gin.H{
        "user":           user,
        "recipe_count":   recipeCount,
        "feedback_count": feedbackCount,
    })
}

// UpdateUserProfile updates user profile information
func UpdateUserProfile(c *gin.Context) {
    id := c.Param("id")
    
    var user models.User
    if err := config.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    var updateData models.User
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
        return
    }
    
    // Don't allow updating sensitive fields through this endpoint
    updateData.Password = ""
    updateData.Username = ""
    updateData.Email = ""
    
    if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile"})
        return
    }
    
    // Remove password from response
    user.Password = ""
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Profile updated successfully",
        "user":    user,
    })
}

// GetAllUsers fetches all users (admin function)
func GetAllUsers(c *gin.Context) {
    var users []models.User
    if err := config.DB.Select("id, username, email, first_name, last_name, bio, avatar_url, is_active, joined_at").Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "users": users,
    })
}

// GetUserRecipes fetches all recipes by a specific user
func GetUserRecipes(c *gin.Context) {
    userID := c.Param("id")
    
    var recipes []models.Recipe
    if err := config.DB.Preload("User").Preload("Feedbacks").Where("user_id = ?", userID).Find(&recipes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user recipes"})
        return
    }
    
    // Calculate average ratings
    for i := range recipes {
        var avgRating float64
        config.DB.Model(&models.Feedback{}).Where("recipe_id = ?", recipes[i].ID).Select("AVG(rating)").Scan(&avgRating)
        recipes[i].AverageRating = avgRating
    }
    
    c.JSON(http.StatusOK, gin.H{
        "recipes": recipes,
    })
}
