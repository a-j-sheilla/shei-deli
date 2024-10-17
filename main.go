package main

import (
    "fmt"
    "log"
    "sheideli/db"
    "sheideli/models"
)

func main() {
    // Connect to the SQLite database
    database, err := db.ConnectSQLite()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer database.Close()

    // Run migrations to create tables if they don't exist
    err = db.RunMigrations(database)
    if err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Example: Insert a recipe
    err = models.InsertRecipe(database, "Vegan Burger", "Delicious plant-based burger with veggies.")
    if err != nil {
        log.Fatal("Failed to insert recipe:", err)
    }

    // Example: Fetch and display all recipes
    recipes, err := models.GetAllRecipes(database)
    if err != nil {
        log.Fatal("Failed to get recipes:", err)
    }

    fmt.Println("Recipes:")
    for _, recipe := range recipes {
        fmt.Printf("ID: %d, Name: %s, Description: %s\n", recipe.ID, recipe.Name, recipe.Description)
    }
}
