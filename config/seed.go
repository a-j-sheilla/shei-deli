package config

import (
    "log"
    "shei-deli/models"
    "golang.org/x/crypto/bcrypt"
)

// SeedDatabase creates initial data for the application
func SeedDatabase() {
    // Create a default admin user
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return
    }

    adminUser := models.User{
        Username:  "admin",
        Email:     "admin@shei-deli.com",
        Password:  string(hashedPassword),
        FirstName: "Admin",
        LastName:  "User",
        Bio:       "Administrator of Shei-deli recipe platform",
        IsActive:  true,
    }

    // Check if admin user already exists
    var existingUser models.User
    result := DB.Where("username = ? OR email = ?", adminUser.Username, adminUser.Email).First(&existingUser)
    if result.Error != nil {
        // User doesn't exist, create it
        if err := DB.Create(&adminUser).Error; err != nil {
            log.Printf("Error creating admin user: %v", err)
            return
        } else {
            log.Println("Admin user created successfully")
        }
    } else {
        log.Println("Admin user already exists")
        adminUser = existingUser // Use existing user
    }

    // Create sample recipes for each category
    sampleRecipes := []models.Recipe{
        {
            Title:        "Quinoa Buddha Bowl",
            Description:  "A nutritious and colorful vegan bowl packed with protein and vegetables",
            Ingredients:  "1 cup quinoa, 1 cup chickpeas, 2 cups mixed vegetables (broccoli, carrots, bell peppers), 1 avocado, 2 tbsp tahini, 1 tbsp lemon juice, salt and pepper to taste",
            Instructions: "1. Cook quinoa according to package instructions. 2. Roast vegetables at 400°F for 20 minutes. 3. Mix tahini with lemon juice for dressing. 4. Assemble bowl with quinoa, vegetables, chickpeas, and avocado. 5. Drizzle with dressing.",
            Category:     models.VeganMeals,
            PrepTime:     15,
            CookTime:     25,
            Servings:     2,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Mini Veggie Pizzas",
            Description:  "Fun and healthy mini pizzas that kids will love",
            Ingredients:  "4 whole wheat English muffins, 1/2 cup pizza sauce, 1 cup mozzarella cheese, 1/2 cup diced vegetables (bell peppers, mushrooms, tomatoes)",
            Instructions: "1. Preheat oven to 375°F. 2. Split English muffins and toast lightly. 3. Spread pizza sauce on each half. 4. Add cheese and vegetables. 5. Bake for 10-12 minutes until cheese melts.",
            Category:     models.KidsMeals,
            PrepTime:     10,
            CookTime:     12,
            Servings:     4,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Grilled Chicken Salad",
            Description:  "Light and protein-rich salad perfect for weight loss",
            Ingredients:  "2 chicken breasts, 4 cups mixed greens, 1 cucumber, 1 cup cherry tomatoes, 1/4 cup balsamic vinegar, 1 tbsp olive oil",
            Instructions: "1. Season and grill chicken breasts until cooked through. 2. Slice chicken and let cool. 3. Mix greens, cucumber, and tomatoes. 4. Top with sliced chicken. 5. Drizzle with balsamic vinegar and olive oil.",
            Category:     models.WeightLossMeals,
            PrepTime:     10,
            CookTime:     15,
            Servings:     2,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Protein Power Smoothie Bowl",
            Description:  "High-calorie, nutrient-dense smoothie bowl for healthy weight gain",
            Ingredients:  "1 banana, 1/2 cup oats, 2 tbsp peanut butter, 1 cup whole milk, 1 scoop protein powder, 1 tbsp honey, granola and nuts for topping",
            Instructions: "1. Blend banana, oats, peanut butter, milk, protein powder, and honey until smooth. 2. Pour into bowl. 3. Top with granola, nuts, and additional fruit as desired.",
            Category:     models.WeightGainMeals,
            PrepTime:     5,
            CookTime:     0,
            Servings:     1,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
    }

    // Create sample recipes if they don't exist
    for _, recipe := range sampleRecipes {
        var existingRecipe models.Recipe
        result := DB.Where("title = ?", recipe.Title).First(&existingRecipe)
        if result.Error != nil {
            // Recipe doesn't exist, create it
            if err := DB.Create(&recipe).Error; err != nil {
                log.Printf("Error creating sample recipe '%s': %v", recipe.Title, err)
            } else {
                log.Printf("Sample recipe '%s' created successfully", recipe.Title)
            }
        }
    }

    log.Println("Database seeding completed")
}
