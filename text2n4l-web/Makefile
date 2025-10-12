.PHONY: all build clean dev-server fmt help lint run test test-api

# Build the application
build:
	@echo "Building text2n4l-web..."
	go build -o bin/text2n4l-web ./cmd/web

# Run the API server for development
dev-server: build
	@echo "Starting API development server..."
	./bin/text2n4l-web

# Run API tests only
test-api:
	@echo "Running API tests..."
	go test -v ./internal/web/handlers_api_test.go

# Run all tests (API only)
test:
	@echo "Running all API tests..."
	go test ./internal/web/handlers_api_test.go

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint the code
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/text2n4l-web
	rm -f coverage.out coverage.html
	rm -f *.prof
	rm -f cpu.prof mem.prof
	rm -rf tmp/

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  dev-server   - Start API development server"
	@echo "  test-api     - Run API tests only"
	@echo "  test         - Run all API tests"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run code linters"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"