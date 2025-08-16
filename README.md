# Shei-deli Recipe Application

A community-driven recipe sharing platform built with Go, Gin, GORM, and SQLite.

## Features

### Recipe Categories
- **Vegan Meals**: Plant-based recipes for vegan lifestyle
- **Kids' Meals**: Fun and healthy recipes that kids will love
- **Weight Loss Meals**: Light, nutritious recipes for weight management
- **Weight Gain Meals**: High-calorie, nutrient-dense recipes for healthy weight gain

### Core Functionality
- **Community-driven**: Users can upload and share their own recipes
- **Feedback System**: Rate and comment on recipes (1-5 star rating)
- **Category Filtering**: Browse recipes by specific categories
- **User Management**: User registration and profile management
- **Recipe Management**: Full CRUD operations for recipes

## Project Structure

```
shei-deli/
├── config/
│   ├── database.go    # Database configuration and initialization
│   └── seed.go        # Initial data seeding
├── controllers/
│   ├── recipe_controllers.go    # Recipe-related endpoints
│   ├── feedback_controllers.go  # Feedback and rating endpoints
│   └── user_controllers.go      # User management endpoints
├── models/
│   ├── recipe.go      # Recipe model with categories
│   ├── feedback.go    # Feedback and rating model
│   └── user.go        # User model
├── routes/
│   └── routes.go      # API route definitions
├── main.go            # Application entry point
├── main_test.go       # Test suite
└── README.md          # This file
```

## API Endpoints

### Health Check
- `GET /health` - Application health status

### Categories
- `GET /api/v1/categories` - Get all available recipe categories

### Recipes
- `GET /api/v1/recipes` - Get all recipes (with optional category filter)
- `GET /api/v1/recipes/:id` - Get specific recipe by ID
- `POST /api/v1/recipes` - Create new recipe
- `PUT /api/v1/recipes/:id` - Update existing recipe
- `DELETE /api/v1/recipes/:id` - Delete recipe
- `GET /api/v1/recipes/category/:category` - Get recipes by category
- `GET /api/v1/recipes/top-rated` - Get top-rated recipes
- `GET /api/v1/recipes/search` - Search recipes (Spoonacular API integration)

### Feedback
- `POST /api/v1/feedback` - Add feedback/rating to recipe
- `GET /api/v1/feedback/recipe/:recipeId` - Get all feedback for a recipe
- `PUT /api/v1/feedback/:id` - Update feedback
- `DELETE /api/v1/feedback/:id` - Delete feedback

### Users
- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/login` - User login
- `GET /api/v1/users` - Get all users (admin)
- `GET /api/v1/users/:id` - Get user profile
- `PUT /api/v1/users/:id` - Update user profile
- `GET /api/v1/users/:id/recipes` - Get user's recipes

## Installation and Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd shei-deli
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Run tests**
   ```bash
   go test -v
   ```

The application will start on `http://localhost:8080`

## Database

The application uses SQLite as the database, which will be automatically created as `shei_deli.db` in the project root when you first run the application.

### Initial Data

The application automatically seeds the database with:
- An admin user (username: `admin`, password: `admin123`)
- Sample recipes for each category

## Usage Examples

### Create a New Recipe
```bash
curl -X POST http://localhost:8080/api/v1/recipes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Quinoa Buddha Bowl",
    "description": "A nutritious vegan bowl",
    "ingredients": "1 cup quinoa, mixed vegetables, tahini",
    "instructions": "Cook quinoa, roast vegetables, assemble bowl",
    "category": "vegan_meals",
    "prep_time": 15,
    "cook_time": 25,
    "servings": 2,
    "difficulty": "Easy",
    "user_id": 1
  }'
```

### Get Recipes by Category
```bash
curl http://localhost:8080/api/v1/recipes/category/vegan_meals
```

### Add Feedback to a Recipe
```bash
curl -X POST http://localhost:8080/api/v1/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "recipe_id": 1,
    "user_id": 1,
    "comment": "Delicious and healthy!",
    "rating": 5
  }'
```

### Register a New User
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "foodlover",
    "email": "foodlover@example.com",
    "password": "password123",
    "first_name": "Food",
    "last_name": "Lover"
  }'
```

## Recipe Categories

### Vegan Meals (`vegan_meals`)
Plant-based recipes that exclude all animal products.

### Kids' Meals (`kids_meals`)
Fun, nutritious, and kid-friendly recipes.

### Weight Loss Meals (`weight_loss_meals`)
Light, low-calorie recipes for healthy weight management.

### Weight Gain Meals (`weight_gain_meals`)
High-calorie, nutrient-dense recipes for healthy weight gain.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is open source and available under the MIT License.
