// Shei-deli Recipe Application JavaScript

// API Base URL
const API_BASE = '/api/v1';

// DOM Content Loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

// Initialize Application
function initializeApp() {
    // Add click handlers for category cards
    const categoryCards = document.querySelectorAll('.category-card');
    categoryCards.forEach(card => {
        card.addEventListener('click', function() {
            const category = this.dataset.category;
            if (category) {
                window.location.href = `/category/${category}`;
            }
        });
    });

    // Add click handlers for recipe cards
    const recipeCards = document.querySelectorAll('.recipe-card');
    recipeCards.forEach(card => {
        card.addEventListener('click', function() {
            const recipeId = this.dataset.recipeId;
            if (recipeId) {
                window.location.href = `/recipe/${recipeId}`;
            }
        });
    });

    // Initialize forms
    initializeForms();
}

// Initialize Forms
function initializeForms() {
    // Recipe form submission
    const recipeForm = document.getElementById('recipeForm');
    if (recipeForm) {
        recipeForm.addEventListener('submit', handleRecipeSubmission);
    }

    // Feedback form submission
    const feedbackForm = document.getElementById('feedbackForm');
    if (feedbackForm) {
        feedbackForm.addEventListener('submit', handleFeedbackSubmission);
    }

    // User registration form
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleUserRegistration);
    }
}

// Handle Recipe Submission
async function handleRecipeSubmission(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const recipeData = {
        title: formData.get('title'),
        description: formData.get('description'),
        ingredients: formData.get('ingredients'),
        instructions: formData.get('instructions'),
        category: formData.get('category'),
        prep_time: parseInt(formData.get('prep_time')) || 0,
        cook_time: parseInt(formData.get('cook_time')) || 0,
        servings: parseInt(formData.get('servings')) || 1,
        difficulty: formData.get('difficulty'),
        user_id: 1 // Default user for now
    };

    try {
        showLoading('Saving recipe...');
        const response = await fetch(`${API_BASE}/recipes`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(recipeData)
        });

        if (response.ok) {
            const result = await response.json();
            showSuccess('Recipe saved successfully!');
            setTimeout(() => {
                window.location.href = `/recipe/${result.ID}`;
            }, 1500);
        } else {
            const error = await response.json();
            showError(error.error || 'Failed to save recipe');
        }
    } catch (error) {
        showError('Network error. Please try again.');
    } finally {
        hideLoading();
    }
}

// Handle Feedback Submission
async function handleFeedbackSubmission(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const feedbackData = {
        recipe_id: parseInt(formData.get('recipe_id')),
        user_id: 1, // Default user for now
        comment: formData.get('comment'),
        rating: parseInt(formData.get('rating'))
    };

    try {
        showLoading('Submitting feedback...');
        const response = await fetch(`${API_BASE}/feedback`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(feedbackData)
        });

        if (response.ok) {
            showSuccess('Feedback submitted successfully!');
            event.target.reset();
            setTimeout(() => {
                location.reload();
            }, 1500);
        } else {
            const error = await response.json();
            showError(error.error || 'Failed to submit feedback');
        }
    } catch (error) {
        showError('Network error. Please try again.');
    } finally {
        hideLoading();
    }
}

// Handle User Registration
async function handleUserRegistration(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const userData = {
        username: formData.get('username'),
        email: formData.get('email'),
        password: formData.get('password'),
        first_name: formData.get('first_name'),
        last_name: formData.get('last_name'),
        bio: formData.get('bio')
    };

    // Validate password confirmation
    const confirmPassword = formData.get('confirm_password');
    if (userData.password !== confirmPassword) {
        showError('Passwords do not match');
        return;
    }

    try {
        showLoading('Creating account...');
        const response = await fetch(`${API_BASE}/users/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });

        if (response.ok) {
            showSuccess('Account created successfully!');
            setTimeout(() => {
                window.location.href = '/';
            }, 1500);
        } else {
            const error = await response.json();
            showError(error.error || 'Failed to create account');
        }
    } catch (error) {
        showError('Network error. Please try again.');
    } finally {
        hideLoading();
    }
}

// Utility Functions
function showLoading(message = 'Loading...') {
    const loadingDiv = document.createElement('div');
    loadingDiv.id = 'loadingMessage';
    loadingDiv.className = 'alert alert-info';
    loadingDiv.innerHTML = `<div class="loading"></div> ${message}`;
    
    const container = document.querySelector('.container');
    container.insertBefore(loadingDiv, container.firstChild);
}

function hideLoading() {
    const loadingDiv = document.getElementById('loadingMessage');
    if (loadingDiv) {
        loadingDiv.remove();
    }
}

function showSuccess(message) {
    showAlert(message, 'success');
}

function showError(message) {
    showAlert(message, 'error');
}

function showAlert(message, type) {
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type}`;
    alertDiv.textContent = message;
    alertDiv.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 1rem 1.5rem;
        border-radius: 5px;
        color: white;
        font-weight: 500;
        z-index: 1000;
        animation: slideIn 0.3s ease;
        background: ${type === 'success' ? '#28a745' : '#dc3545'};
    `;
    
    document.body.appendChild(alertDiv);
    
    setTimeout(() => {
        alertDiv.remove();
    }, 3000);
}

// Star Rating Component
function createStarRating(rating, interactive = false) {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
        const star = document.createElement('span');
        star.className = 'star';
        star.innerHTML = i <= rating ? '★' : '☆';
        star.style.color = i <= rating ? '#ffc107' : '#ddd';
        
        if (interactive) {
            star.style.cursor = 'pointer';
            star.addEventListener('click', () => setRating(i));
            star.addEventListener('mouseover', () => highlightStars(i));
        }
        
        stars.push(star);
    }
    return stars;
}

function setRating(rating) {
    const ratingInput = document.getElementById('rating');
    if (ratingInput) {
        ratingInput.value = rating;
    }
    updateStarDisplay(rating);
}

function highlightStars(rating) {
    const stars = document.querySelectorAll('.star');
    stars.forEach((star, index) => {
        star.style.color = index < rating ? '#ffc107' : '#ddd';
    });
}

function updateStarDisplay(rating) {
    const stars = document.querySelectorAll('.star');
    stars.forEach((star, index) => {
        star.innerHTML = index < rating ? '★' : '☆';
        star.style.color = index < rating ? '#ffc107' : '#ddd';
    });
}

// Search functionality
function searchRecipes(query) {
    if (query.length < 2) return;
    
    fetch(`${API_BASE}/recipes?search=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(data => {
            displaySearchResults(data.recipes);
        })
        .catch(error => {
            console.error('Search error:', error);
        });
}

function displaySearchResults(recipes) {
    const resultsContainer = document.getElementById('searchResults');
    if (!resultsContainer) return;
    
    resultsContainer.innerHTML = '';
    
    if (recipes.length === 0) {
        resultsContainer.innerHTML = '<p>No recipes found.</p>';
        return;
    }
    
    recipes.forEach(recipe => {
        const recipeElement = createRecipeCard(recipe);
        resultsContainer.appendChild(recipeElement);
    });
}

function createRecipeCard(recipe) {
    const card = document.createElement('div');
    card.className = 'recipe-card';
    card.dataset.recipeId = recipe.ID;
    
    card.innerHTML = `
        <div class="recipe-image"></div>
        <div class="recipe-content">
            <h3 class="recipe-title">${recipe.title}</h3>
            <p class="recipe-description">${recipe.description}</p>
            <div class="recipe-meta">
                <span>${recipe.prep_time + recipe.cook_time} min</span>
                <div class="rating">
                    <span class="stars">${'★'.repeat(Math.floor(recipe.average_rating))}</span>
                    <span>${recipe.average_rating.toFixed(1)}</span>
                </div>
            </div>
        </div>
    `;
    
    card.addEventListener('click', () => {
        window.location.href = `/recipe/${recipe.ID}`;
    });
    
    return card;
}
