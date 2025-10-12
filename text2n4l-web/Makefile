.PHONY: all benchmark build clean dev install-deps profile run test test-all test-coverage test-verbose

# Makefile for text2n4l-web development and testing

# Default target
all: test build

# Build the application
build:
	@echo "Building text2n4l-web..."
	go build -o bin/text2n4l-web ./cmd/web

# Install development dependencies
install-deps:
	@echo "Installing development dependencies..."
	go install github.com/cosmtrek/air@latest
	go mod tidy

# Run the application in development mode with hot reload
dev:
	@echo "Starting development server with hot reload..."
	air

# Run the application
run: build
	@echo "Starting text2n4l-web server..."
	./bin/text2n4l-web

# Run all tests
test:
	@echo "Running all tests..."
	go test ./...

# Run all tests including templates
test-all:
	@echo "Running all tests and template validation..."
	go test ./...
	$(MAKE) test-templates

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Test templates specifically
test-templates:
	@echo "Running template validation tests..."
	go test -v ./internal/web -run TestTemplate

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage analysis..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with coverage and open report
test-coverage-view: test-coverage
	@if command -v xdg-open >/dev/null 2>&1; then \
		xdg-open coverage.html; \
	elif command -v open >/dev/null 2>&1; then \
		open coverage.html; \
	else \
		echo "Open coverage.html in your browser"; \
	fi

# Run specific package tests
test-analyzer:
	@echo "Running analyzer tests..."
	go test -v ./internal/analyzer

test-web:
	@echo "Running web handler tests..."
	go test -v ./internal/web

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Run benchmarks for specific packages
benchmark-analyzer:
	@echo "Running analyzer benchmarks..."
	go test -bench=. -benchmem ./internal/analyzer

benchmark-web:
	@echo "Running web benchmarks..."
	go test -bench=. -benchmem ./internal/web

# Profile CPU usage
profile-cpu:
	@echo "Starting server for CPU profiling..."
	@echo "Run your load tests, then visit http://localhost:8080/debug/pprof/profile"
	@echo "Or run: go tool pprof http://localhost:8080/debug/pprof/profile"
	make run

# Profile memory usage
profile-memory:
	@echo "Starting server for memory profiling..."
	@echo "Run your load tests, then visit http://localhost:8080/debug/pprof/heap"
	@echo "Or run: go tool pprof http://localhost:8080/debug/pprof/heap"
	make run

# Update golden files for regression testing
update-golden:
	@echo "Updating golden files..."
	UPDATE_GOLDEN=1 go test ./internal/analyzer -run TestGoldenFiles

# Lint the code
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/text2n4l-web
	rm -f coverage.out coverage.html
	rm -f *.prof
	rm -f cpu.prof mem.prof
	rm -rf tmp/

# Full development setup
setup: install-deps
	@echo "Creating directories..."
	mkdir -p bin tmp static templates
	@echo "Running template validation..."
	make test-templates
	@echo "Setup complete! Run 'make dev' to start development server."

# Template validation (run before commits)
validate-templates: test-templates
	@echo "✅ All templates validated successfully"

# CI/CD targets
ci-test: test-race test-coverage test-templates lint
	@echo "CI tests completed"

# Performance testing with load
load-test:
	@echo "Running load tests (requires server to be running)..."
	@echo "Install ab (Apache Bench) if not available"
	@if command -v ab >/dev/null 2>&1; then \
		echo "Testing homepage..."; \
		ab -n 1000 -c 10 http://localhost:8080/; \
		echo "Testing convert endpoint..."; \
		ab -n 100 -c 5 -p testdata/sample.txt -T 'application/x-www-form-urlencoded' http://localhost:8080/convert; \
	else \
		echo "Apache Bench (ab) not found. Install with: sudo apt-get install apache2-utils"; \
	fi

# Generate test data
generate-testdata:
	@echo "Generating test data..."
	@echo "This is a simple test document." > internal/analyzer/testdata/simple.txt
	@echo "It contains several sentences for testing the N4L conversion process." >> internal/analyzer/testdata/simple.txt
	@echo "Each sentence represents a potential unit of analysis." >> internal/analyzer/testdata/simple.txt

# Docker targets (if needed)
docker-build:
	@echo "Building Docker image..."
	docker build -t text2n4l-web .

docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 text2n4l-web

# Help target
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  dev             - Start development server with hot reload"
	@echo "  run             - Run the application"
	@echo "  test            - Run all tests"
	@echo "  test-verbose    - Run tests with verbose output"
	@echo "  test-templates  - Run template validation tests"
	@echo "  test-race       - Run tests with race detection"
	@echo "  test-coverage   - Run tests with coverage analysis"
	@echo "  validate-templates - Validate all template syntax"
	@echo "  benchmark       - Run benchmarks"
	@echo "  profile-cpu     - Start server for CPU profiling"
	@echo "  profile-memory  - Start server for memory profiling"
	@echo "  update-golden   - Update golden files for regression testing"
	@echo "  lint            - Run code linters"
	@echo "  fmt             - Format code"
	@echo "  clean           - Clean build artifacts"
	@echo "  setup           - Full development setup"
	@echo "  ci-test         - Run CI tests (race, coverage, templates, lint)"
	@echo "  load-test       - Run load tests"
	@echo "  help            - Show this help message"