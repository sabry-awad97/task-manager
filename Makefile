# Binary name
BINARY_NAME=task-manager.exe

# Build directory
BUILD_DIR=build

# Main package path
MAIN_PACKAGE=./cmd/task-manager

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Build flags
BUILD_FLAGS=-v -ldflags="-s -w"

.PHONY: all build run clean test lint vet deps help

all: clean build

build:
	@echo "Building..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@$(GOBUILD) $(BUILD_FLAGS) -o "$(BUILD_DIR)/$(BINARY_NAME)" $(MAIN_PACKAGE)

run: build
	@echo "Running..."
	@"$(BUILD_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	@if exist "$(BUILD_DIR)" rd /s /q "$(BUILD_DIR)"
	@$(GOCLEAN)

test:
	@echo "Running tests..."
	@$(GOTEST) ./... -v

lint:
	@echo "Running linter..."
	@$(GOLINT) run

vet:
	@echo "Running vet..."
	@$(GOVET) ./...

deps:
	@echo "Installing dependencies..."
	@$(GOGET) -v -t -d ./...

help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run     - Build and run the application"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make test    - Run tests"
	@echo "  make lint    - Run linter"
	@echo "  make vet     - Run go vet"
	@echo "  make deps    - Install dependencies"
	@echo "  make help    - Show this help message"
