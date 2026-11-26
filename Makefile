GO ?= go

.PHONY: help run test test-verbose test-race fmt vet tidy build check

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-14s %s\\n", $$1, $$2}'

run: ## Run the health-check API
	$(GO) run ./cmd/main.go

test: ## Run all tests with verbose output
	$(GO) test -v ./...

test-race: ## Run tests with the race detector (requires CGO and a C compiler)
	$(GO) test -race ./...

fmt: ## Format Go source files
	$(GO) fmt ./...

vet: ## Run Go static analysis
	$(GO) vet ./...

tidy: ## Clean up go.mod and go.sum
	$(GO) mod tidy

build: ## Build all packages
	$(GO) build ./...

check: fmt vet test build ## Run the main quality checks