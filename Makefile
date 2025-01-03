# Variables
APP_NAME=asana
BUILD_DIR=build
VERSION := $(shell git describe --tags --abbrev=0 --always)
LDFLAGS := -ldflags "-X github.com/timwehrle/asana/pkg/version.Version=${VERSION}"

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

.PHONY: build
build: ## Run build
	@echo "Running build"
	$(GOBUILD) ${LDFLAGS} -o $(BUILD_DIR)/$(APP_NAME)

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -coverprofile=c.out ./...
	$(GOCMD) tool cover -html=c.out

.PHONY: lint
lint: ## Run linter
	@echo "Running lint..."
	$(LINT)

.PHONY: fmt
fmt: ## Run formatter
	@echo "Formatting code..."
	$(GOFMT) -w -s .

.PHONY: audit
audit: ## Audit code
	@echo "Running audit..."
	$(GOMOD) tidy
	$(GOMOD) verify
	$(GOVET) ./...
	$(GORUN) $(VULN) ./...

.PHONY: release
release: ## Run GoReleaser
	@echo "Releasing..."
	goreleaser release --clean

.PHONY: help
help: ## Show available commands
	@grep -E '^[a-zA-Z_/.-]+:.*?##' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = "##"}; {printf "%-20s %s\n", $$1, $$2}'
