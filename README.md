# Shei-deli Recipe Application

A community-driven recipe sharing platform built with Go, Gin, GORM, and SQLite with a beautiful web interface.

## Features

### Recipe Categories
- **Plant-Based Meals**: Vegan/vegetarian options (no animal products)
- **Kids' Meals**: Fun, simple, and nutritious meals for children
- **Light Meals (Weight Loss)**: Low-calorie, balanced recipes
- **Hearty Meals (Weight Gain)**: High-calorie, energy-packed recipes
- **Meat Stews**: Beef, chicken, goat, lamb, and other meat-based stews
- **Veggie Stews**: Lentil, bean, mushroom, and vegetable stews
- **Seafood & Fish Stews**: Fish stews, seafood mixes, and ocean-inspired flavors
- **Fusion Stews**: Cultural and traditional varieties (e.g., goulash, curries)
- **Soups**: Warm, comforting soups
- **Drinks**: Smoothies, juices, teas, and other beverages
- **Pastries**: Baked goods such as cakes, cookies, pies, and breads

### Core Functionality
- **Beautiful Web Interface**: Responsive design with category images and intuitive navigation
- **Community-driven**: Users can upload and share their own recipes
- **Feedback System**: Interactive 5-star rating and comment system
- **Category Filtering**: Browse recipes by specific categories with visual category cards
- **User Management**: User registration and profile management
- **Recipe Management**: Full CRUD operations for recipes
- **Mobile-Friendly**: Responsive design that works on all devices

## Project Structure

```
shei-deli/
├── config/
│   ├── database.go         # Database configuration and initialization
│   ├── seed.go            # Initial data seeding
│   └── template_helpers.go # Template helper functions
├── controllers/
│   ├── recipe_controllers.go    # Recipe-related API endpoints
│   ├── feedback_controllers.go  # Feedback and rating API endpoints
│   ├── user_controllers.go      # User management API endpoints
│   └── web_controllers.go       # Web interface controllers
├── models/
│   ├── recipe.go      # Recipe model with 10 categories
│   ├── feedback.go    # Feedback and rating model
│   └── user.go        # User model
├── routes/
│   └── routes.go      # API and web route definitions
├── static/
│   ├── css/style.css  # Application styles
│   └── js/app.js      # Interactive JavaScript
├── templates/
│   ├── base.html      # Base template layout
│   ├── index.html     # Home page with category grid
│   ├── category.html  # Category recipe listings
│   ├── recipe.html    # Recipe detail page
│   ├── add-recipe.html # Recipe creation form
│   └── register.html  # User registration form
├── images/            # Category images
├── main.go            # Application entry point
├── main_test.go       # Test suite
└── README.md          # This file
```

## Web Interface

The application features a beautiful, responsive web interface accessible at `http://localhost:8080`

### Web Pages
- `/` - Home page with category grid and featured recipes
- `/category/:category` - Category-specific recipe listings
- `/recipe/:id` - Detailed recipe view with ratings and comments
- `/add-recipe` - Recipe creation form
- `/register` - User registration
- `/recipes` - All recipes with pagination
- `/about` - About page

### Features
- **Visual Category Navigation**: Click on category images to browse recipes
- **Interactive Recipe Cards**: Hover effects and click-to-view functionality
- **Star Rating System**: Interactive 5-star rating for recipes
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile
- **Form Validation**: Client-side validation for better user experience

## API Endpoints

### Health Check
- `GET /health` - Application health status

### Categories
- `GET /api/v1/categories` - Get all 11 available recipe categories

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
    "category": "plant_based_meals",
    "prep_time": 15,
    "cook_time": 25,
    "servings": 2,
    "difficulty": "Easy",
    "user_id": 1
  }'
```

### Get Recipes by Category
```bash
curl http://localhost:8080/api/v1/recipes/category/plant_based_meals
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

### Plant-Based Meals (`plant_based_meals`)
Vegan/vegetarian options with no animal products.

### Kids' Meals (`kids_meals`)
Fun, simple, and nutritious meals for children.

### Light Meals (`light_meals`)
Low-calorie, balanced recipes for weight management.

### Hearty Meals (`hearty_meals`)
High-calorie, energy-packed recipes for healthy weight gain.

### Meat Stews (`meat_stews`)
Beef, chicken, goat, lamb, and other meat-based stews.

### Veggie Stews (`veggie_stews`)
Lentil, bean, mushroom, and vegetable stews.

### Seafood & Fish Stews (`seafood_stews`)
Fish stews, seafood mixes, and ocean-inspired flavors.

### Fusion Stews (`fusion_stews`)
Cultural and traditional varieties (e.g., goulash, curries).

### Soups (`soups`)
Warm, comforting soups for every occasion.

### Drinks (`drinks`)
Smoothies, juices, teas, and other beverages.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is open source and available under the MIT License.
