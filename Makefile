.PHONY: help dev start stop build clean logs test

# Default target
help:
	@echo "Construction Estimation & Bidding Automation Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  make dev     - Start all services in development mode"
	@echo "  make start   - Start all services"
	@echo "  make stop    - Stop all services"
	@echo "  make build   - Build all Docker images"
	@echo "  make clean   - Remove containers and volumes"
	@echo "  make logs    - Tail logs from all services"
	@echo "  make test    - Run tests for all services"

# Start all services in development mode
dev:
	@echo "Starting all services in development mode..."
	docker compose up --build

# Start all services (detached)
start:
	@echo "Starting all services..."
	docker compose up -d

# Stop all services
stop:
	@echo "Stopping all services..."
	docker compose down

# Build all Docker images
build:
	@echo "Building all Docker images..."
	docker compose build

# Clean up containers, volumes, and images
clean:
	@echo "Cleaning up containers, volumes, and images..."
	docker compose down -v
	docker system prune -f

# Tail logs from all services
logs:
	@echo "Tailing logs from all services..."
	docker compose logs -f

# Run tests for all services
test:
	@echo "Running tests for all services..."
	@echo "\n=== Testing Backend (Go) ==="
	cd backend && go test -v ./...
	@echo "\n=== Testing AI Service (Python) ==="
	cd ai_service && pip install -q -r requirements.txt -r requirements-dev.txt && pytest -v
	@echo "\n=== Testing App (React Native) ==="
	cd app && npm install --silent && npm test
	@echo "\n=== All tests completed ==="
