# 🚀 Shei-deli Startup Guide

This guide will help you get the Shei-deli recipe application running with both backend and frontend.

## 📋 Prerequisites

- Go 1.22+ installed
- Git installed
- Web browser

## 🛠️ Setup Instructions

### 1. Fix Go Module Dependencies

Due to module verification issues, run the following commands:

```bash
# Set environment variables to bypass verification
export GOPROXY=direct
export GOSUMDB=off

# Clean and download dependencies
go clean -modcache
go mod tidy
```

### 2. Start the Backend Server

```bash
# Run the application
GOPROXY=direct GOSUMDB=off go run main.go
```

You should see output like:
```
Database connected and migrated!
Admin user created successfully
Sample recipe 'Quinoa Buddha Bowl' created successfully
...
Starting Shei-deli server on :8080...
Web interface: http://localhost:8080
API documentation: http://localhost:8080/api/v1/categories
```

### 3. Test the Integration

Once the server is running, you can test it in several ways:

#### Option A: Web Interface
Open your browser and go to: `http://localhost:8080`

#### Option B: Integration Test Page
Open the integration test file: `file:///path/to/shei-deli/integration_test.html`

#### Option C: Command Line Test
Run the test script:
```bash
./test_backend.sh
```

## 🌐 Available Endpoints

### Web Interface
- `http://localhost:8080/` - Home page with category grid
- `http://localhost:8080/category/plant_based_meals` - Plant-based recipes
- `http://localhost:8080/add-recipe` - Add new recipe form
- `http://localhost:8080/register` - User registration

### API Endpoints
- `GET /health` - Health check
- `GET /api/v1/categories` - All 11 recipe categories
- `GET /api/v1/recipes` - All recipes
- `GET /api/v1/recipes/category/{category}` - Recipes by category
- `POST /api/v1/recipes` - Create new recipe
- `POST /api/v1/feedback` - Add recipe feedback
- `POST /api/v1/users/register` - Register new user

## 🎯 Testing the Features

### 1. Browse Categories
- Visit the home page
- Click on any category card (Plant-Based, Kids' Meals, etc.)
- View recipes in that category

### 2. Add a Recipe
- Go to `/add-recipe`
- Fill out the form with recipe details
- Select a category
- Submit the form

### 3. Rate Recipes
- Click on any recipe to view details
- Use the star rating system
- Add comments

### 4. Register Users
- Go to `/register`
- Create a new user account
- Login and manage recipes

## 🗂️ Database

The application uses SQLite with the database file `shei_deli.db` created automatically.

### Sample Data
The application seeds the database with:
- Admin user (username: `admin`, password: `admin123`)
- Sample recipes for all 11 categories
- Initial feedback and ratings

## 🎨 Category Images

The application uses the provided category images:
- `vegan.jpeg` - Plant-Based Meals
- `kids-meals.jpeg` - Kids' Meals
- `light-meals.jpeg` - Light Meals
- `hearty-meals.jpeg` - Hearty Meals
- `stews.jpeg` - Meat Stews
- `vegetable-stews.jpeg` - Veggie Stews
- `fish&sea-food.jpeg` - Seafood Stews
- `fusion.jpeg` - Fusion Stews
- `soups.jpeg` - Soups
- `drinks&smoothies.jpeg` - Drinks
- `pastries.jpeg` - Pastries

## 🔧 Troubleshooting

### Server Won't Start
1. Check Go version: `go version`
2. Clear module cache: `go clean -modcache`
3. Use bypass flags: `GOPROXY=direct GOSUMDB=off go run main.go`

### Database Issues
1. Delete `shei_deli.db` file
2. Restart the server to recreate with fresh data

### Port Already in Use
1. Check what's using port 8080: `lsof -i :8080`
2. Kill the process or change the port in `main.go`

### Images Not Loading
1. Ensure images are in the `images/` directory
2. Check file permissions
3. Verify image file names match the controller mappings

## 📱 Mobile Testing

The interface is responsive and works on mobile devices:
- Test on different screen sizes
- Verify touch interactions work
- Check category grid layout

## 🚀 Production Deployment

For production deployment:
1. Build the binary: `go build -o shei-deli .`
2. Set environment variables for database
3. Use a reverse proxy (nginx) for static files
4. Enable HTTPS
5. Set up proper logging

## 📊 Monitoring

Monitor the application:
- Health endpoint: `/health`
- Database file size
- Response times
- Error logs

## 🎉 Success!

If everything is working, you should be able to:
- ✅ Browse beautiful category cards with background images
- ✅ View recipes filtered by category
- ✅ Add new recipes through the web form
- ✅ Rate and comment on recipes
- ✅ Register new users
- ✅ Use all API endpoints

Your Shei-deli recipe sharing platform is now ready for the community! 🍽️
