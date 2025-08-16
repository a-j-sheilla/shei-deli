package main

import (
    "html/template"
    "log"
    "shei-deli/config"
    "shei-deli/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize database connection and run migrations
    config.InitDatabase()

    // Seed the database with initial data
    config.SeedDatabase()

    // Setup Gin with template functions
    gin.SetMode(gin.ReleaseMode) // Set to debug mode for development

    // Setup routes
    router := routes.SetupRoutes()

    // Set template functions
    router.SetFuncMap(config.GetTemplateFunctions())

    // Reload templates (for development)
    router.LoadHTMLGlob("templates/*")

    // Start the server
    log.Println("Starting Shei-deli server on :8080...")
    log.Println("Web interface: http://localhost:8080")
    log.Println("API documentation: http://localhost:8080/api/v1/categories")

    if err := router.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
