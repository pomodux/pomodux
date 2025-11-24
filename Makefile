.PHONY: build clean test install help

# Version information
VERSION ?= 0.1.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -X main.version=$(VERSION) \
           -X main.buildTime=$(BUILD_TIME) \
           -X main.gitCommit=$(GIT_COMMIT)

# Binary names
POMODUX_BIN := pomodux
STATS_BIN := pomodux-stats

# Build directory
BUILD_DIR := bin

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: build-pomodux build-stats ## Build both binaries

build-pomodux: ## Build pomodux binary
	@echo "Building $(POMODUX_BIN)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(POMODUX_BIN) ./cmd/pomodux

build-stats: ## Build pomodux-stats binary
	@echo "Building $(STATS_BIN)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(STATS_BIN) ./cmd/pomodux-stats

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

install: build ## Install both binaries to GOPATH/bin
	@echo "Installing binaries..."
	@go install -ldflags "$(LDFLAGS)" ./cmd/pomodux
	@go install -ldflags "$(LDFLAGS)" ./cmd/pomodux-stats
	@echo "Installed $(POMODUX_BIN) and $(STATS_BIN)"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint (if installed)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

mod-tidy: ## Tidy go.mod dependencies
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify

