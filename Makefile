# Variables
APP_NAME=alfie
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
	$(GOTEST) -v -coverprofile=c.out ./...
	$(GOCMD) tool cover -html=c.out

.PHONY: test/cover
test/cover:
	go test -v -coverprofile=c.out ./...
	go tool cover -html=c.out

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
