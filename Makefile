.PHONY: build test clean install

BINARY_NAME=cf-ecslogs-plugin
BUILD_DIR=bin
GO=go
GOFMT=gofmt

# Get the current git commit and version
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION=1.0.0

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)"

all: build

# Build the plugin
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .

# Check code formatting
fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$($(GOFMT) -l .)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

# Vet code
vet:
	@echo "Vetting code..."
	$(GO) vet ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install the plugin locally
install: build
	@echo "Installing plugin..."
	cf install-plugin -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Plugin installed successfully"

# Uninstall the plugin
uninstall:
	@echo "Uninstalling plugin..."
	cf uninstall-plugin ECSLogsPlugin || true
	@echo "Plugin uninstalled"

# Reinstall the plugin (uninstall + install)
reinstall: uninstall install

# Run linting (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@echo "Running linter..."
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the plugin for current platform"
	@echo "  build-all      - Build for all platforms (Linux, macOS, Windows)"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  fmt-check      - Check if code is formatted"
	@echo "  vet            - Run go vet"
	@echo "  deps           - Install/update dependencies"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Build and install plugin to CF CLI"
	@echo "  uninstall      - Uninstall plugin from CF CLI"
	@echo "  reinstall      - Uninstall and install plugin"
	@echo "  lint           - Run golangci-lint"
	@echo "  help           - Show this help message"
