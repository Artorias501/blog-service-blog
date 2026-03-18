# Makefile for blog-service
# Provides build, run, test, and clean targets

# Binary name
BINARY_NAME=blog-service
BINARY_UNIX=$(BINARY_NAME)_unix

# Build directory
BUILD_DIR=./bin

# Main package
MAIN_PACKAGE=./cmd/server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build run test clean help deps fmt lint

# Default target
all: deps build

# ============================================================================
# BUILD TARGETS
# ============================================================================

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

## build-linux: Build the binary for Linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_PACKAGE)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_UNIX)"

# ============================================================================
# RUN TARGETS
# ============================================================================

## run: Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

## run-dev: Run the application in development mode with hot reload (requires air)
run-dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

# ============================================================================
# TEST TARGETS
# ============================================================================

## test: Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## test-short: Run short tests
test-short:
	@echo "Running short tests..."
	$(GOTEST) -v -short ./...

# ============================================================================
# CLEAN TARGETS
# ============================================================================

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# ============================================================================
# DEPENDENCY TARGETS
# ============================================================================

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify
	@echo "Dependencies downloaded"

## deps-tidy: Tidy dependencies
deps-tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	@echo "Dependencies tidied"

# ============================================================================
# CODE QUALITY TARGETS
# ============================================================================

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Code formatted"

## lint: Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "Running linter..."
	golangci-lint run ./...
	@echo "Linting complete"

# ============================================================================
# DATABASE TARGETS
# ============================================================================

## db-migrate: Run database migrations
db-migrate:
	@echo "Running database migrations..."
	@echo "Migrations are run automatically on startup"

## db-reset: Reset database
db-reset:
	@echo "Resetting database..."
	rm -f blog.db
	@echo "Database reset complete"

# ============================================================================
# HELP TARGET
# ============================================================================

## help: Show this help message
help:
	@echo "Blog Service - Available targets:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Examples:"
	@echo "  make build        # Build the binary"
	@echo "  make run          # Build and run the application"
	@echo "  make test         # Run all tests"
	@echo "  make clean        # Clean build artifacts"
