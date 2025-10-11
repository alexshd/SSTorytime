.PHONY: build clean fmt help run test vet

# Makefile for N4L Text Converter Web Application

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build both web and CLI applications"
	@echo "  build-web - Build only web application"
	@echo "  build-cli - Build only CLI application"
	@echo "  run      - Run the development server"
	@echo "  demo     - Run CLI demo with test file"
	@echo "  test     - Run all tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  fmt      - Format Go code"
	@echo "  vet      - Run go vet"
	@echo "  deps     - Download dependencies"

# Build the application
build:
	@echo "Building N4L Web Converter..."
	go build -o bin/n4l-web cmd/web/main.go
	@echo "Building N4L CLI Converter..."
	go build -o bin/n4l-cli cmd/cli/main.go
	@echo "Build complete: bin/n4l-web, bin/n4l-cli"

# Build only web version
build-web:
	@echo "Building N4L Web Converter..."
	go build -o bin/n4l-web cmd/web/main.go
	@echo "Build complete: bin/n4l-web"

# Build only CLI version
build-cli:
	@echo "Building N4L CLI Converter..."
	go build -o bin/n4l-cli cmd/cli/main.go
	@echo "Build complete: bin/n4l-cli"

# Run the development server
run:
	@echo "Starting N4L Web Converter server..."
	go run cmd/web/main.go

# Run CLI version with demo file
demo:
	@echo "Running CLI demo..."
	go run cmd/cli/main.go testdata/demo.txt

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/n4l-web bin/n4l-cli
	rm -f coverage.out coverage.html
	go clean

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Development setup
dev-setup: deps fmt vet test
	@echo "Development environment ready!"

# Docker build (optional)
docker-build:
	@echo "Building Docker image..."
	docker build -t n4l-web .

# Install for system-wide use
install: build
	@echo "Installing n4l-web to /usr/local/bin..."
	sudo cp bin/n4l-web /usr/local/bin/
	@echo "Installation complete!"

# Check for security issues
security:
	@echo "Running security checks..."
	go list -json -deps ./... | jq -r '.Module | select(.Replace == null) | .Path + " " + .Version' | sort | uniq

# Performance benchmark
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# All checks (for CI/CD)
check: fmt vet test
	@echo "All checks passed!"