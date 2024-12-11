# Variables
APP_NAME=jodot
BUILD_DIR=

# Commands
GOCMD := go 
GOTEST := $(GOCMD) test
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOCLEAN := $(GOCMD) clean
LINT := golangci-lint run
GOFMT := gofmt

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Run linter
.PHONY: lint
lint:
	@echo "Running lint..."
	$(LINT)

# Run formatter
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w -s .
