.PHONY: help build run test clean docker-build docker-run docker-stop swagger

# Default target
help:
	@echo "iSHARE Task API - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose services"
	@echo ""
	@echo "Documentation:"
	@echo "  swagger      - Generate Swagger documentation"
	@echo ""
	@echo "Database:"
	@echo "  db-create    - Create PostgreSQL database"
	@echo "  db-migrate   - Run database migrations"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

# Run the application locally
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t ishare-task-api .

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go

# Create PostgreSQL database (requires psql)
db-create:
	@echo "Creating PostgreSQL database..."
	createdb ishare_tasks || echo "Database already exists or psql not available"

# Run database migrations (this happens automatically on startup)
db-migrate:
	@echo "Database migrations run automatically on application startup"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Run with hot reload (requires air)
dev:
	@echo "Running with hot reload..."
	air

# Show logs
logs:
	@echo "Showing application logs..."
	docker-compose logs -f app

# Database logs
db-logs:
	@echo "Showing database logs..."
	docker-compose logs -f postgres

# Reset database
db-reset:
	@echo "Resetting database..."
	docker-compose down -v
	docker-compose up -d postgres
	sleep 5
	docker-compose up -d app

# Full setup (install deps, generate docs, build)
setup: deps swagger build
	@echo "Setup complete!"

# Production build
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go 