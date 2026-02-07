.PHONY: build run test lint clean help

# Build the application
build:
	@echo "Building backend..."
	go build -o bin/api ./cmd/api

# Build for production (optimized)
build-release:
	@echo "Building for production..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api ./cmd/api

# Run the application
run:
	@echo "Running backend..."
	go run ./cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the application with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Starting development server with hot reload..."
	air

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-release  - Build for production"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy dependencies"
	@echo "  clean          - Clean build artifacts"
	@echo "  install-tools  - Install development tools"
	@echo "  dev            - Run with hot reload (requires air)"
