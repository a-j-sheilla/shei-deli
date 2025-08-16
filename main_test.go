package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "shei-deli/config"
    "shei-deli/models"
    "shei-deli/routes"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB() {
    // Use in-memory SQLite for testing
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to test database")
    }
    
    config.DB = db
    
    // Auto-migrate the schema
    err = db.AutoMigrate(&models.Recipe{}, &models.Feedback{}, &models.User{})
    if err != nil {
        panic("Failed to migrate test database")
    }
}

func TestHealthEndpoint(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := routes.SetupRoutes()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/health", nil)
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("Failed to parse response: %v", err)
    }
    
    if response["status"] != "healthy" {
        t.Errorf("Expected status 'healthy', got %v", response["status"])
    }
}

func TestGetCategories(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := routes.SetupRoutes()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("Failed to parse response: %v", err)
    }
    
    categories, ok := response["categories"].([]interface{})
    if !ok {
        t.Errorf("Expected categories array in response")
    }
    
    if len(categories) != 4 {
        t.Errorf("Expected 4 categories, got %d", len(categories))
    }
}

func TestCreateUser(t *testing.T) {
    setupTestDB()
    gin.SetMode(gin.TestMode)
    router := routes.SetupRoutes()
    
    user := map[string]interface{}{
        "username":   "testuser",
        "email":      "test@example.com",
        "password":   "password123",
        "first_name": "Test",
        "last_name":  "User",
    }
    
    jsonData, _ := json.Marshal(user)
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusCreated {
        t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
    }
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("Failed to parse response: %v", err)
    }
    
    if response["message"] != "User created successfully" {
        t.Errorf("Expected success message, got %v", response["message"])
    }
}

func TestCreateRecipe(t *testing.T) {
    setupTestDB()
    gin.SetMode(gin.TestMode)
    router := routes.SetupRoutes()
    
    // First create a user
    user := models.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    config.DB.Create(&user)
    
    recipe := map[string]interface{}{
        "title":        "Test Recipe",
        "description":  "A test recipe",
        "ingredients":  "Test ingredients",
        "instructions": "Test instructions",
        "category":     "vegan_meals",
        "prep_time":    15,
        "cook_time":    30,
        "servings":     4,
        "difficulty":   "Easy",
        "user_id":      user.ID,
    }
    
    jsonData, _ := json.Marshal(recipe)
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/recipes", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusCreated {
        t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
    }
    
    var response models.Recipe
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("Failed to parse response: %v", err)
    }
    
    if response.Title != "Test Recipe" {
        t.Errorf("Expected title 'Test Recipe', got %s", response.Title)
    }
}

func TestGetRecipes(t *testing.T) {
    setupTestDB()
    gin.SetMode(gin.TestMode)
    router := routes.SetupRoutes()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/recipes", nil)
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("Failed to parse response: %v", err)
    }
    
    if _, ok := response["recipes"]; !ok {
        t.Errorf("Expected 'recipes' field in response")
    }
}
