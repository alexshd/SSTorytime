.PHONY: all build clean dev-server fmt help lint profile-cpu profile-heap run test test-api

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
	go test -v ./internal/web/

# Run all tests (API only)
test:
	@echo "Running all API tests..."
	go test ./internal/web/

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
	rm -f cpu.prof mem.prof heap.prof
	rm -rf tmp/

# Profile CPU usage (requires server to be running)
profile-cpu:
	@echo "Capturing CPU profile (30 seconds)..."
	@echo "Make sure the server is running with 'make dev-server'"
	curl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof
	@echo "CPU profile saved to cpu.prof"
	@echo "Analyze with: go tool pprof cpu.prof"

# Profile heap/memory usage (requires server to be running)
profile-heap:
	@echo "Capturing heap profile..."
	@echo "Make sure the server is running with 'make dev-server'"
	curl http://localhost:8080/debug/pprof/heap -o heap.prof
	@echo "Heap profile saved to heap.prof"
	@echo "Analyze with: go tool pprof heap.prof"

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  dev-server   - Start API development server"
	@echo "  test-api     - Run API tests only"
	@echo "  test         - Run all API tests"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run code linters"
	@echo "  profile-cpu  - Capture CPU profile (server must be running)"
	@echo "  profile-heap - Capture heap/memory profile (server must be running)"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"