.PHONY: help dev start stop build clean logs test prod-build prod-start prod-stop prod-logs validate

# Default target
help:
	@echo "Construction Estimation & Bidding Automation Platform"
	@echo ""
	@echo "Development Commands:"
	@echo "  make validate      - Validate development setup"
	@echo "  make dev           - Start all services in development mode"
	@echo "  make start         - Start all services"
	@echo "  make stop          - Stop all services"
	@echo "  make build         - Build all Docker images"
	@echo "  make clean         - Remove containers and volumes"
	@echo "  make logs          - Tail logs from all services"
	@echo "  make test          - Run tests for all services"
	@echo ""
	@echo "Production Commands:"
	@echo "  make prod-build    - Build production Docker images"
	@echo "  make prod-start    - Start production services"
	@echo "  make prod-stop     - Stop production services"
	@echo "  make prod-logs     - Tail production logs"
	@echo "  make push-images   - Build and push images to registry"
	@echo ""
	@echo "Quality & Setup:"
	@echo "  make validate      - Validate development environment setup"
	@echo "  make setup-hooks   - Install pre-commit hooks"

# Validate development setup
validate:
	@echo "Validating development setup..."
	@./scripts/validate-setup.sh

# Setup pre-commit hooks
setup-hooks:
	@echo "Installing pre-commit hooks..."
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit install; \
		echo "Pre-commit hooks installed successfully!"; \
	else \
		echo "Error: pre-commit not found. Install with: pip install pre-commit"; \
		exit 1; \
	fi

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

# Build production Docker images
prod-build:
	@echo "Building production Docker images..."
	docker build -f backend/Dockerfile.production -t construction-backend:latest ./backend
	docker build -f ai_service/Dockerfile.production -t construction-ai-service:latest ./ai_service
	docker build -f app/Dockerfile.production -t construction-frontend:latest ./app
	@echo "Production images built successfully!"

# Start production services
prod-start:
	@echo "Starting production services..."
	@if [ ! -f .env.production ]; then \
		echo "Error: .env.production not found. Copy .env.production.example and configure it."; \
		exit 1; \
	fi
	docker-compose -f docker-compose.production.yml --env-file .env.production up -d
	@echo "Production services started!"

# Stop production services
prod-stop:
	@echo "Stopping production services..."
	docker-compose -f docker-compose.production.yml down

# Tail production logs
prod-logs:
	@echo "Tailing production logs..."
	docker-compose -f docker-compose.production.yml logs -f

# Build and push images to container registry
push-images:
	@echo "Building and pushing images to container registry..."
	./scripts/push-images.sh
