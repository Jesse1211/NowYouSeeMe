.PHONY: help backend frontend demo db-setup db-reset db-migrate db-connect install install-backend install-frontend install-sdk build build-backend build-frontend docs

# Default target
.DEFAULT_GOAL := help

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

##@ Help

help: ## Show this help message
	@echo "$(CYAN)NowYouSeeMe - Event Sourcing Platform$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make $(CYAN)<target>$(NC)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(CYAN)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

backend: ## Start backend server
	@echo "$(GREEN)Starting backend server...$(NC)"
	cd backend && go run main.go

frontend: ## Start frontend dev server
	@echo "$(GREEN)Starting frontend dev server...$(NC)"
	cd frontend && npm run dev

demo: ## Generate quick demo data (6 agents)
	@echo "$(GREEN)Generating quick demo data...$(NC)"
	cd sdk && python3 scripts/seed_database.py --preset quick

demo-full: ## Generate full dataset (17 agents with rich history)
	@echo "$(GREEN)Generating full dataset...$(NC)"
	cd sdk && python3 scripts/seed_database.py --preset full

demo-mbti: ## Generate MBTI diversity (32 agents, one per type)
	@echo "$(GREEN)Generating MBTI diversity dataset...$(NC)"
	cd sdk && python3 scripts/seed_database.py --preset mbti

demo-custom: ## Generate custom dataset (use NUM=N ENTRIES=E)
	@echo "$(GREEN)Generating custom dataset...$(NC)"
	cd sdk && python3 scripts/seed_database.py --custom -n $(NUM) -e $(ENTRIES)

##@ Database

# Load environment variables from .env file if it exists
-include .env
export

db-setup: ## Setup PostgreSQL database
	@echo "$(GREEN)Setting up PostgreSQL database...$(NC)"
	@bash backend/scripts/setup_db.sh

db-reset: db-setup ## Reset database (drop + recreate)
	@echo "$(GREEN)Database reset complete$(NC)"

db-migrate: ## Run migrations only
	@echo "$(GREEN)Running database migrations...$(NC)"
	@psql -h $${DB_HOST:-localhost} -U $${DB_USER:-liuzhenhua} -d $${DB_NAME:-nowyouseeme} -f backend/migrations/001_create_event_sourcing_schema.sql

db-connect: ## Connect to database using psql
	@echo "$(GREEN)Connecting to database...$(NC)"
	@psql -h $${DB_HOST:-localhost} -U $${DB_USER:-liuzhenhua} -d $${DB_NAME:-nowyouseeme}

##@ Setup

install-backend: ## Install backend dependencies
	@echo "$(GREEN)Installing backend dependencies...$(NC)"
	cd backend && go mod download

install-frontend: ## Install frontend dependencies
	@echo "$(GREEN)Installing frontend dependencies...$(NC)"
	cd frontend && npm install

install-sdk: ## Install SDK
	@echo "$(GREEN)Installing SDK...$(NC)"
	cd sdk && pip3 install -e .

install: ## Install all dependencies
	@echo "$(GREEN)Installing all dependencies...$(NC)"
	@make install-backend
	@make install-frontend
	@make install-sdk

##@ Build

build-backend: ## Build backend binary
	@echo "$(GREEN)Building backend...$(NC)"
	cd backend && go build -o nowyouseeme

build-frontend: ## Build frontend for production
	@echo "$(GREEN)Building frontend...$(NC)"
	cd frontend && npm run build

build: ## Build all
	@make build-backend
	@make build-frontend

##@ Documentation

docs: ## Show documentation overview
	@echo "$(CYAN)Documentation:$(NC)"
	@echo "  README.md                           - Project overview"
	@echo "  docs/superpowers/specs/*.md         - Design specifications"
	@echo "  sdk/examples/generate_sample_data.py - Sample data generator"
	@echo ""
	@echo "$(CYAN)Quick Start:$(NC)"
	@echo "  1. make install                     - Install all dependencies"
	@echo "  2. make db-setup                    - Setup PostgreSQL database"
	@echo "  3. make backend                     - Start backend (Terminal 1)"
	@echo "  4. make frontend                    - Start frontend (Terminal 2)"
	@echo "  5. make demo                        - Generate sample data (Terminal 3)"
	@echo "  6. Visit http://localhost:3000      - View the gallery"
