# Variables
APP_NAME=asana
BUILD_DIR=

# Commands
GOCMD := go
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet
GOTEST := $(GOCMD) test
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOCLEAN := $(GOCMD) clean
LINT := golangci-lint run
GOFMT := gofmt
VULN := golang.org/x/vuln/cmd/govulncheck@latest

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

# Audit code
.PHONY: audit
audit:
	@echo "Running audit..."
	$(GOMOD) tidy
	$(GOMOD) verify
	$(GOVET) ./...
	$(GORUN) $(VULN) ./...

# Show available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make test       Run tests and generate coverage report"
	@echo "  make lint       Run linter"
	@echo "  make fmt        Format code"
	@echo "  make audit      Audit code for issues"
