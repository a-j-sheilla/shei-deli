# Shei-deli Recipe Application Makefile

.PHONY: help build run test clean demo deps

# Default target
help:
	@echo "Shei-deli Recipe Application"
	@echo "============================"
	@echo ""
	@echo "Available commands:"
	@echo "  make deps    - Download and install dependencies"
	@echo "  make build   - Build the application"
	@echo "  make run     - Run the application"
	@echo "  make test    - Run tests"
	@echo "  make demo    - Run the demo script (requires running app)"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make help    - Show this help message"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Build the application
build: deps
	@echo "🔨 Building application..."
	go build -o shei-deli .

# Run the application
run: deps
	@echo "🚀 Starting Shei-deli server..."
	go run main.go

# Run tests
test: deps
	@echo "🧪 Running tests..."
	go test -v ./...

# Run demo (requires the app to be running)
demo:
	@echo "🎬 Running demo..."
	@if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then \
		echo "❌ Application is not running. Please start it with 'make run' first."; \
		exit 1; \
	fi
	./demo.sh

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f shei-deli
	rm -f shei_deli.db
	go clean

# Development workflow
dev: clean build run

# Show application info
info:
	@echo "Shei-deli Recipe Application"
	@echo "============================"
	@echo "A community-driven recipe sharing platform"
	@echo ""
	@echo "Features:"
	@echo "- 🥗 Vegan Meals"
	@echo "- 👶 Kids' Meals" 
	@echo "- 📉 Weight Loss Meals"
	@echo "- 📈 Weight Gain Meals"
	@echo "- 👥 User Management"
	@echo "- ⭐ Rating & Feedback System"
	@echo "- 🔍 Category Filtering"
	@echo ""
	@echo "API Documentation: http://localhost:8080/health"
