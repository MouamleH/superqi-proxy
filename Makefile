# Makefile for superqi-proxy

APP_NAME := superqi-proxy
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := builds
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all build clean test run dev build-all help

# Default target
all: build

# Build for current platform
build:
	@echo "🔨 Building $(APP_NAME) for current platform..."
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .
	@echo "✅ Built $(APP_NAME)"

# Build for all platforms
build-all:
	@echo "🚀 Building $(APP_NAME) for all platforms..."
	./build.sh

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(APP_NAME)
	@echo "✅ Clean completed"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run the application in development mode
run:
	@echo "🏃 Running $(APP_NAME) in development mode..."
	go run .

# Development build with race detection
dev:
	@echo "🔧 Building $(APP_NAME) for development..."
	go build -race -ldflags "$(LDFLAGS)" -o $(APP_NAME)-dev .
	@echo "✅ Built $(APP_NAME)-dev with race detection"

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build for current platform"
	@echo "  build-all - Build for all platforms using build.sh"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  run       - Run in development mode"
	@echo "  dev       - Build with race detection"
	@echo "  help      - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  VERSION   - Set build version (default: git describe or 'dev')"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make build-all"
	@echo "  VERSION=1.0.0 make build-all"
