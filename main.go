package main

import (
    "shei-deli/config"
    "shei-deli/routes"
)

func main() {
    // Initialize the database
    config.InitDatabase()

    // Set up routes
    router := routes.SetupRoutes()
    router.Run(":8080")
}
