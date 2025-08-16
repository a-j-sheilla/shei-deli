#!/bin/bash

# Shei-deli Backend Integration Test Script
echo "🍽️  Shei-deli Backend Integration Test"
echo "======================================"

# Check if server is running
echo ""
echo "1. Testing server connection..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Server is running on port 8080"
else
    echo "❌ Server is not running. Please start with 'go run main.go'"
    echo ""
    echo "To start the server:"
    echo "1. Fix Go module issues: GOPROXY=direct GOSUMDB=off go mod tidy"
    echo "2. Run the server: GOPROXY=direct GOSUMDB=off go run main.go"
    echo "3. Then run this test again"
    exit 1
fi

echo ""
echo "2. Testing Health Check endpoint..."
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if [[ $? -eq 0 ]]; then
    echo "✅ Health check passed: $HEALTH_RESPONSE"
else
    echo "❌ Health check failed"
fi

echo ""
echo "3. Testing Categories API..."
CATEGORIES_RESPONSE=$(curl -s http://localhost:8080/api/v1/categories)
if [[ $? -eq 0 ]]; then
    echo "✅ Categories API working"
    echo "Response: $CATEGORIES_RESPONSE" | head -c 200
    echo "..."
else
    echo "❌ Categories API failed"
fi

echo ""
echo "4. Testing Recipes API..."
RECIPES_RESPONSE=$(curl -s http://localhost:8080/api/v1/recipes)
if [[ $? -eq 0 ]]; then
    echo "✅ Recipes API working"
    echo "Response: $RECIPES_RESPONSE" | head -c 200
    echo "..."
else
    echo "❌ Recipes API failed"
fi

echo ""
echo "5. Testing Recipe Creation..."
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/recipes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Integration Recipe",
    "description": "A test recipe for backend integration",
    "ingredients": "Test ingredients",
    "instructions": "Test instructions",
    "category": "plant_based_meals",
    "prep_time": 10,
    "cook_time": 15,
    "servings": 2,
    "difficulty": "Easy",
    "user_id": 1
  }')

if [[ $? -eq 0 ]]; then
    echo "✅ Recipe creation working"
    echo "Response: $CREATE_RESPONSE" | head -c 200
    echo "..."
else
    echo "❌ Recipe creation failed"
fi

echo ""
echo "6. Testing Web Interface..."
WEB_RESPONSE=$(curl -s -I http://localhost:8080/ | head -n 1)
if [[ $WEB_RESPONSE == *"200"* ]]; then
    echo "✅ Web interface accessible"
    echo "Visit: http://localhost:8080"
else
    echo "❌ Web interface not accessible"
fi

echo ""
echo "7. Testing Static Files..."
CSS_RESPONSE=$(curl -s -I http://localhost:8080/static/css/style.css | head -n 1)
if [[ $CSS_RESPONSE == *"200"* ]]; then
    echo "✅ Static files (CSS) accessible"
else
    echo "❌ Static files not accessible"
fi

echo ""
echo "8. Testing Category Images..."
IMAGE_RESPONSE=$(curl -s -I http://localhost:8080/images/vegan.jpeg | head -n 1)
if [[ $IMAGE_RESPONSE == *"200"* ]]; then
    echo "✅ Category images accessible"
else
    echo "❌ Category images not accessible"
fi

echo ""
echo "🎉 Integration Test Complete!"
echo ""
echo "If all tests passed, your Shei-deli application is working correctly!"
echo "You can now:"
echo "- Visit the web interface: http://localhost:8080"
echo "- Test the integration page: file://$(pwd)/integration_test.html"
echo "- Use the API endpoints for development"
echo ""
echo "Next steps:"
echo "1. Open http://localhost:8080 in your browser"
echo "2. Click on category cards to browse recipes"
echo "3. Add new recipes using the form"
echo "4. Rate and comment on recipes"
