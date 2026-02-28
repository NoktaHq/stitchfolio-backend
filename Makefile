.PHONY: setup help build run run-dev run-qa launch launch-qa serve migrate migrate-qa migrate-dev wire swagger pre-commit install-hooks

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

# Go tools (used by setup and pre-commit)
WIRE_VERSION ?= latest
SWAG_VERSION ?= v1.16.6

help:
	@echo "Available commands:"
	@echo "  make setup         - Set up development environment (Go, wire, swag)"
	@echo "  make pre-commit    - Run wire, swagger, build (use before commit)"
	@echo "  make install-hooks - Install git pre-commit hook (runs pre-commit on every commit)"
	@echo ""
	@echo "  make build         - Build binary to app/app"
	@echo "  make wire         - Regenerate wire_gen.go"
	@echo "  make swagger      - Regenerate Swagger docs"
	@echo ""
	@echo "  make run           - Run with prod config"
	@echo "  make run-dev       - Run with dev config"
	@echo "  make run-qa        - Run with QA config"
	@echo "  make launch        - Run built binary (prod)"
	@echo "  make launch-qa     - Run built binary (QA)"
	@echo "  make serve         - Build and launch"
	@echo ""
	@echo "  make migrate       - Run migrations (prod)"
	@echo "  make migrate-dev   - Run migrations (dev)"
	@echo "  make migrate-qa    - Run migrations (QA)"

setup:
	@echo "$(YELLOW)Starting development environment setup...$(NC)"
	@echo ""

	# Check Go
	@if command -v go > /dev/null 2>&1; then \
		echo "$(GREEN)✓ Go is already installed (version: $$(go version | cut -d' ' -f3))$(NC)"; \
	else \
		echo "$(RED)Go not found. Please install Go: https://go.dev/doc/install$(NC)"; \
		exit 1; \
	fi
	@echo ""

	# Download Go modules
	@echo "$(YELLOW)Downloading Go modules...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Go modules downloaded$(NC)"
	@echo ""

	# Install wire
	@echo "$(YELLOW)Installing wire...$(NC)"
	@if command -v wire > /dev/null 2>&1; then \
		echo "$(GREEN)✓ wire is already installed$$(wire 2>/dev/null | head -1)$(NC)"; \
	else \
		go install github.com/google/wire/cmd/wire@$(WIRE_VERSION) && \
		echo "$(GREEN)✓ wire installed successfully$(NC)" || (echo "$(RED)Failed to install wire$(NC)"; exit 1); \
	fi
	@echo ""

	# Install swag (swagger)
	@echo "$(YELLOW)Installing swag (Swagger)...$(NC)"
	@if command -v swag > /dev/null 2>&1; then \
		echo "$(GREEN)✓ swag is already installed (version: $$(swag --version 2>/dev/null || swag version 2>/dev/null))$(NC)"; \
	else \
		go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION) && \
		echo "$(GREEN)✓ swag installed successfully$(NC)" || (echo "$(RED)Failed to install swag$(NC)"; exit 1); \
	fi
	@echo ""

	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)✓ Dev setup is done!$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo ""
	@echo "Optional: run 'make install-hooks' to run wire + swagger + build before every commit."

# Pre-commit: regenerate code and verify build. Run this before committing.
pre-commit: wire swagger build
	@echo "$(GREEN)✓ pre-commit checks passed$(NC)"

# Use git pre-commit hook so every commit runs wire, swagger, build
install-hooks:
	@chmod +x .githooks/pre-commit 2>/dev/null || true
	@echo "$(GREEN)✓ Pre-commit hook is at .githooks/pre-commit$(NC)"
	@echo ""
	@echo "Enable it (run once): git config core.hooksPath .githooks"
	@echo "Then every 'git commit' will run: make wire, make swagger, make build"
	@echo "and auto-stage generated files (wire_gen.go, docs/)."

# --- Build ---
build:
	@echo "$(YELLOW)Building...$(NC)"
	go build -tags netgo -ldflags '-s -w' -o app/app
	@echo "$(GREEN)✓ Build successful$(NC)"

# --- Run ---
run:
	go run main.go --configFile config/prod.yaml

run-dev:
	go run main.go --configFile config/dev.yaml

run-qa:
	go run main.go --configFile config/qa.yaml

launch:
	./app/app --configFile=config/prod.yaml

launch-qa:
	./app/app --configFile=config/qa.yaml

serve: build launch

# --- Migrate ---
migrate:
	go run main.go --configFile config/prod.yaml --migrate true

migrate-qa:
	go run main.go --configFile config/qa.yaml --migrate true

migrate-dev:
	go run main.go --configFile config/dev.yaml --migrate true

# --- Dev tools ---
wire:
	@echo "$(YELLOW)Regenerating wire_gen.go...$(NC)"
	@cd internal/di && go generate ./...
	@echo "$(GREEN)✓ wire done$(NC)"

swagger:
	@echo "$(YELLOW)Regenerating Swagger docs...$(NC)"
	swag fmt
	swag init --parseDependency --parseInternal
	go run ./scripts/patch_swagger_security.go
	@echo "$(GREEN)✓ swagger done$(NC)"
