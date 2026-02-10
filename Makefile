.PHONY: help start stop test crud create read update delete clean backend frontend demo

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
	@echo "$(CYAN)NowYouSeeMe - Makefile Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make $(CYAN)<target>$(NC)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(CYAN)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

backend: ## Start backend server
	@echo "$(GREEN)Starting backend server...$(NC)"
	cd backend && go run main.go

frontend: ## Start frontend dev server
	@echo "$(GREEN)Starting frontend dev server...$(NC)"
	cd frontend && npm run dev

demo: ## Generate sample demo data
	@echo "$(GREEN)Generating demo data...$(NC)"
	cd sdk && python3 examples/generate_sample_data.py

##@ Testing

test: ## Run all tests
	@echo "$(GREEN)Running all tests...$(NC)"
	@cd sdk && python3 tests/test_create.py
	@cd sdk && python3 tests/test_read.py
	@cd sdk && python3 tests/test_update.py
	@cd sdk && python3 tests/test_delete.py

test-full: ## Run full CRUD test suite
	@echo "$(GREEN)Running full CRUD test suite...$(NC)"
	cd sdk && python3 tests/test_crud_full.py

test-create: ## Test CREATE only
	cd sdk && python3 tests/test_create.py

test-read: ## Test READ only
	cd sdk && python3 tests/test_read.py

test-update: ## Test UPDATE only
	cd sdk && python3 tests/test_update.py

test-delete: ## Test DELETE only
	cd sdk && python3 tests/test_delete.py

##@ CRUD Operations

crud: ## Execute full CRUD cycle (auto-generated)
	@echo "$(GREEN)Executing CRUD cycle...$(NC)"
	cd sdk && python3 scripts/crud.py

crud-custom: ## Execute CRUD with custom name (usage: make crud-custom NAME="MyAgent")
	@echo "$(GREEN)Executing CRUD with custom name: $(NAME)$(NC)"
	cd sdk && python3 scripts/crud.py "$(NAME)"

create: ## Create visualization (auto-generated)
	@echo "$(GREEN)Creating visualization...$(NC)"
	cd sdk && python3 scripts/create.py

create-custom: ## Create with custom name (usage: make create-custom NAME="MyAgent")
	cd sdk && python3 scripts/create.py "$(NAME)"

read: ## Read all visualizations
	@echo "$(GREEN)Reading all visualizations...$(NC)"
	cd sdk && python3 scripts/read.py

read-one: ## Read one visualization (usage: make read-one ID=<viz_id>)
	cd sdk && python3 scripts/read.py $(ID)

update: ## Update visualization (usage: make update ID=<viz_id>)
	cd sdk && python3 scripts/update.py $(ID)

update-custom: ## Update with custom name (usage: make update-custom ID=<viz_id> NAME="NewName")
	cd sdk && python3 scripts/update.py $(ID) "$(NAME)"

delete: ## Delete visualization (usage: make delete ID=<viz_id>)
	cd sdk && python3 scripts/delete.py $(ID) --force

##@ Utilities

random: ## Add random visualization
	@echo "$(GREEN)Adding random visualization...$(NC)"
	cd sdk && python3 scripts/add_random.py

random-many: ## Add multiple random visualizations (usage: make random-many N=10)
	@echo "$(GREEN)Adding $(N) random visualizations...$(NC)"
	@for i in $$(seq 1 $(N)); do \
		cd sdk && python3 scripts/add_random.py; \
		sleep 0.3; \
	done

list: ## List all visualizations
	cd sdk && python3 scripts/list_all.py

clean: ## Clear all visualizations
	@echo "$(RED)Clearing all data...$(NC)"
	cd sdk && python3 scripts/clear_all.py

##@ Quick Workflows

quick-test: ## Quick test: clean + crud + read
	@echo "$(GREEN)Quick test workflow...$(NC)"
	@make clean
	@sleep 1
	@make crud
	@sleep 1
	@make read

populate: ## Populate with 10 random visualizations
	@echo "$(GREEN)Populating with 10 random visualizations...$(NC)"
	@make random-many N=10

full-demo: ## Full demo: clean + populate + read
	@echo "$(GREEN)Running full demo...$(NC)"
	@make clean
	@sleep 1
	@make populate
	@sleep 1
	@make list
	@echo ""
	@echo "$(GREEN)✓ Demo complete! View at: http://localhost:3000$(NC)"

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
	@echo "  README.md              - Project overview"
	@echo "  QUICKSTART.md          - Quick start guide"
	@echo "  sdk/QUICK_REFERENCE.md - SDK quick reference"
	@echo "  sdk/SCRIPTS_GUIDE.md   - Scripts usage guide"
	@echo "  sdk/TESTING_GUIDE.md   - Testing guide"
	@echo "  docs/.context/         - Detailed documentation"

##@ Aliases (Short commands)

c: create ## Alias for 'create'
r: read ## Alias for 'read'
u: ## Update (usage: make u ID=<viz_id>)
	@make update ID=$(ID)
d: ## Delete (usage: make d ID=<viz_id>)
	@make delete ID=$(ID)
t: test-full ## Alias for 'test-full'
