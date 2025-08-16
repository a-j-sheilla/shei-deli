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
            Category:     models.PlantBasedMeals,
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
            Category:     models.LightMeals,
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
            Category:     models.HeartyMeals,
            PrepTime:     5,
            CookTime:     0,
            Servings:     1,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Beef and Vegetable Stew",
            Description:  "Hearty beef stew with tender vegetables and rich broth",
            Ingredients:  "2 lbs beef chuck, 4 carrots, 3 potatoes, 2 onions, 3 celery stalks, 4 cups beef broth, 2 tbsp tomato paste, herbs and spices",
            Instructions: "1. Brown beef in large pot. 2. Add onions and cook until soft. 3. Add tomato paste and cook 1 minute. 4. Add broth and bring to boil. 5. Simmer 1.5 hours. 6. Add vegetables and cook 30 minutes more.",
            Category:     models.MeatStews,
            PrepTime:     20,
            CookTime:     120,
            Servings:     6,
            Difficulty:   "Medium",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Lentil and Mushroom Stew",
            Description:  "Rich and satisfying vegetarian stew with lentils and mushrooms",
            Ingredients:  "2 cups green lentils, 1 lb mixed mushrooms, 2 onions, 4 carrots, 4 cups vegetable broth, 2 tbsp olive oil, herbs and spices",
            Instructions: "1. Heat oil in large pot. 2. Sauté onions until soft. 3. Add mushrooms and cook until browned. 4. Add lentils and broth. 5. Simmer 45 minutes until lentils are tender. 6. Add carrots in last 15 minutes.",
            Category:     models.VeggieStews,
            PrepTime:     15,
            CookTime:     45,
            Servings:     4,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Mediterranean Fish Stew",
            Description:  "Fresh seafood stew with Mediterranean flavors",
            Ingredients:  "1 lb white fish, 1/2 lb shrimp, 2 cups diced tomatoes, 1 onion, 3 cloves garlic, 2 cups fish stock, olive oil, herbs",
            Instructions: "1. Heat oil in large pot. 2. Sauté onion and garlic. 3. Add tomatoes and stock. 4. Simmer 15 minutes. 5. Add fish and cook 5 minutes. 6. Add shrimp and cook 3 minutes more.",
            Category:     models.SeafoodStews,
            PrepTime:     15,
            CookTime:     25,
            Servings:     4,
            Difficulty:   "Medium",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Thai Curry Stew",
            Description:  "Fusion stew with Thai curry flavors and coconut milk",
            Ingredients:  "1 lb chicken, 1 can coconut milk, 2 tbsp red curry paste, mixed vegetables, fish sauce, lime juice, fresh herbs",
            Instructions: "1. Cook curry paste in pot until fragrant. 2. Add coconut milk and bring to simmer. 3. Add chicken and cook 15 minutes. 4. Add vegetables and cook until tender. 5. Season with fish sauce and lime juice.",
            Category:     models.FusionStews,
            PrepTime:     10,
            CookTime:     25,
            Servings:     4,
            Difficulty:   "Medium",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Classic Chicken Soup",
            Description:  "Comforting homemade chicken soup with vegetables",
            Ingredients:  "1 whole chicken, 2 carrots, 2 celery stalks, 1 onion, egg noodles, parsley, salt and pepper",
            Instructions: "1. Simmer chicken in water for 1 hour. 2. Remove chicken and shred meat. 3. Strain broth. 4. Add vegetables to broth and cook 15 minutes. 5. Add noodles and chicken, cook until noodles are tender.",
            Category:     models.Soups,
            PrepTime:     15,
            CookTime:     75,
            Servings:     6,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Green Smoothie",
            Description:  "Healthy green smoothie packed with nutrients",
            Ingredients:  "2 cups spinach, 1 banana, 1 apple, 1 cup coconut water, 1 tbsp chia seeds, 1 tbsp honey",
            Instructions: "1. Add all ingredients to blender. 2. Blend until smooth. 3. Add more coconut water if needed for desired consistency. 4. Serve immediately over ice.",
            Category:     models.Drinks,
            PrepTime:     5,
            CookTime:     0,
            Servings:     2,
            Difficulty:   "Easy",
            UserID:       adminUser.ID,
        },
        {
            Title:        "Classic Chocolate Chip Cookies",
            Description:  "Soft and chewy chocolate chip cookies that everyone loves",
            Ingredients:  "2 1/4 cups flour, 1 tsp baking soda, 1 tsp salt, 1 cup butter, 3/4 cup brown sugar, 3/4 cup white sugar, 2 eggs, 2 tsp vanilla, 2 cups chocolate chips",
            Instructions: "1. Preheat oven to 375°F. 2. Mix flour, baking soda, and salt in bowl. 3. Cream butter and sugars. 4. Beat in eggs and vanilla. 5. Gradually add flour mixture. 6. Stir in chocolate chips. 7. Drop spoonfuls on baking sheet. 8. Bake 9-11 minutes.",
            Category:     models.Pastries,
            PrepTime:     15,
            CookTime:     10,
            Servings:     24,
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
