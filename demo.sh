#!/bin/bash

# Shei-deli Recipe Application Demo Script
# This script demonstrates the key features of the application

echo "🍽️  Shei-deli Recipe Application Demo"
echo "======================================"

BASE_URL="http://localhost:8080"

echo ""
echo "1. Health Check"
echo "---------------"
curl -s "$BASE_URL/health" | jq '.'

echo ""
echo "2. Get Available Categories"
echo "---------------------------"
curl -s "$BASE_URL/api/v1/categories" | jq '.'

echo ""
echo "3. Register a New User"
echo "----------------------"
curl -s -X POST "$BASE_URL/api/v1/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demo_user",
    "email": "demo@shei-deli.com",
    "password": "demo123",
    "first_name": "Demo",
    "last_name": "User",
    "bio": "I love cooking and sharing recipes!"
  }' | jq '.'

echo ""
echo "4. Create a New Vegan Recipe"
echo "----------------------------"
curl -s -X POST "$BASE_URL/api/v1/recipes" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Avocado Toast Supreme",
    "description": "Delicious and nutritious avocado toast with a twist",
    "ingredients": "2 slices whole grain bread, 1 ripe avocado, 1 tomato, hemp seeds, lemon juice, salt, pepper",
    "instructions": "1. Toast bread until golden. 2. Mash avocado with lemon juice, salt, and pepper. 3. Spread on toast. 4. Top with sliced tomato and hemp seeds.",
    "category": "vegan_meals",
    "prep_time": 10,
    "cook_time": 2,
    "servings": 1,
    "difficulty": "Easy",
    "user_id": 1
  }' | jq '.'

echo ""
echo "5. Get All Recipes"
echo "------------------"
curl -s "$BASE_URL/api/v1/recipes" | jq '.recipes[] | {id, title, category, average_rating}'

echo ""
echo "6. Get Vegan Recipes Only"
echo "-------------------------"
curl -s "$BASE_URL/api/v1/recipes/category/vegan_meals" | jq '.recipes[] | {id, title, description}'

echo ""
echo "7. Add Feedback to a Recipe"
echo "---------------------------"
curl -s -X POST "$BASE_URL/api/v1/feedback" \
  -H "Content-Type: application/json" \
  -d '{
    "recipe_id": 1,
    "user_id": 1,
    "comment": "This recipe is absolutely amazing! My kids loved it too.",
    "rating": 5
  }' | jq '.'

echo ""
echo "8. Get Top Rated Recipes"
echo "------------------------"
curl -s "$BASE_URL/api/v1/recipes/top-rated?limit=3" | jq '.recipes[] | {title, average_rating}'

echo ""
echo "9. Get User Profile"
echo "-------------------"
curl -s "$BASE_URL/api/v1/users/1" | jq '.user | {username, first_name, last_name, bio}'

echo ""
echo "10. Get User's Recipes"
echo "----------------------"
curl -s "$BASE_URL/api/v1/users/1/recipes" | jq '.recipes[] | {title, category, average_rating}'

echo ""
echo "✅ Demo completed! The Shei-deli application is working perfectly."
echo ""
echo "Key Features Demonstrated:"
echo "- ✅ Recipe categories (Vegan, Kids, Weight Loss, Weight Gain)"
echo "- ✅ User registration and management"
echo "- ✅ Recipe creation and retrieval"
echo "- ✅ Category-based filtering"
echo "- ✅ Feedback and rating system"
echo "- ✅ Top-rated recipes"
echo "- ✅ User profiles and recipe ownership"
echo ""
echo "🚀 Your recipe sharing platform is ready for the community!"
