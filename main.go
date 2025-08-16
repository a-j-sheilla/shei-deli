package main

import (
    "log"
    "shei-deli/config"
    "shei-deli/routes"
)

func main() {
    // Initialize database connection and run migrations
    config.InitDatabase()

    // Seed the database with initial data
    config.SeedDatabase()

    // Setup routes
    router := routes.SetupRoutes()

    // Start the server
    log.Println("Starting Shei-deli server on :8080...")
    if err := router.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
