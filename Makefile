.PHONY: code-review run docs migrate rollback create-migration force-version open-coverage

# Default DB configuration
DB_MASTER_HOST ?= localhost
DB_MASTER_PORT ?= 3306
DB_MASTER_USER ?= user
DB_MASTER_PASSWORD ?= password
DB_MASTER_NAME ?= appointment_management

# Construct DB_URL using the environment variables or fallback to default values
DB_URL := mysql://$(DB_MASTER_USER):$(DB_MASTER_PASSWORD)@tcp($(DB_MASTER_HOST):$(DB_MASTER_PORT))/$(DB_MASTER_NAME)

# Paths
MIGRATIONS_DIR=./migrations

# Tools
MIGRATE_BIN=migrate
SWAG_BIN=swag

open-coverage: 
ifeq ($(shell uname), Darwin) # macOS
	open coverage.html
else ifeq ($(shell uname), Linux)
	xdg-open coverage.html
else ifeq ($(OS), Windows_NT)
	start coverage.html
else
	@echo "Unsupported OS. Please open $(COVERAGE_HTML) manually."
endif

code-review:
	go clean -testcache
	go test -race ./tests/... -coverpkg=./internal/... -coverprofile=coverage.out && \
	go tool cover -func=coverage.out > coverage.txt && \
	go tool cover -html=coverage.out -o coverage.html

	@output=$$(nestif --min 4 ./internal/...); \
	if [ -n "$$output" ]; then \
		echo "$$output"; \
		echo "Error: Detected nested if statements with complexity greater than 3."; \
		exit 1; \
	fi
	
	go run ./tools/code_review/

run:
	go run ./cmd/app/

docs:
	$(SWAG_BIN) init -g cmd/app/main.go

migrate: ## Run all pending migrations
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

rollback: ## Rollback the last migration
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

create-migration: ## Create a new migration file
ifndef NAME
	$(error NAME is required. Usage: make create-migration NAME=your_migration_name)
endif
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

force-version: ## Force the migration version
ifndef VERSION
	$(error VERSION is required. Usage: make force-version VERSION=your_version_number)
endif
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)
