APP := kovago
CMD := ./cmd/kovago
BIN := ./bin/${APP}

GO := go
COMPOSE := docker compose

CONFIG ?= config/local.toml

MIGRATIONS := ./internal/migrator/migrations
GOOSE := $(GO) tool goose

DB_SERVICE := postgres

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available commands
	@echo "KovaGo development commands"
	@echo
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.## / {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run KovaGo
	@echo "→ running KovaGo"
	@KOVAGO_CONFIG=$(CONFIG) $(GO) run $(CMD)

.PHONY: build
build: ## Build KovaGo binary
	@echo "→ bulding $(BIN)"
	@mkdir -p ./bin
	@$(GO) build -trimpath -o $(BIN) $(CMD)

.PHONY: fmt
fmt: ## Format Go source
	@echo "→ formatting"
	@$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "→ vetting"
	@$(GO) vet ./...

.PHONY: test
test: ## Run tests
	@echo "→ testing"
	@$(GO) test ./...

.PHONY: test-race
test-race: ## Run tests with race detector
	@echo "→ testing with race detector"
	@$(GO) test -race ./...

.PHONY: coverage
coverage: ## Generate test coverage report
	@echo "→ generating coverage"
	@$(GO) test -coverprofile=coverage.out ./...
	@$(GO) tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html: coverage ## Open HTML coverage report
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "→ coverage report: coverage.html"

.PHONY: tidy
tidy: ## Tidy Go modules
	@echo "→ tidying modules"
	@$(GO) mod tidy

.PHONY: download
download: ## Dowload Go modules
	@echo "→ downloading modules"
	@$(GO) mod download

.PHONY: check
check: fmt vet test migrations-validate ## Run local quality checks
	@echo "✓ all checks passed"

.PHONY: ci
ci: vet test-race ## Run CI checks
	@echo "✓ CI checks passed"

.PHONY: clean
clean: ## Remove generated files
	@echo "→ cleaning"
	@rm -rf ./bin
	@rm -f coverage.out coverage.html

.PHONY: logs
logs: ## Follow KovaGo log file
	@tail -f logs/kovago.log

.PHONY: db-up
db-up: ## Start PostgreSQL
	@echo "→ starting PostgreSQL"
	@$(COMPOSE) up -d $(DB_SERVICE)
	@echo "→ waiting for PostgreSQL"
	@until $(COMPOSE) exec -T $(DB_SERVICE) pg_isready -U kovago -d kovago >/dev/null 2>&1; do sleep 1; done
	@echo "✓ PostgreSQL is ready"

.PHONY: db-down
db-down: ## Stop PostgreSQL
	@echo "→ stopping PostgreSQL"
	@$(COMPOSE) down

.PHONY: db-status
db-status: ## Show PostgreSQL status
	@$(COMPOSE) ps $(DB_SERVICE)

.PHONY: db-logs
db-logs: ## Follow PostgreSQL logs
	@$(COMPOSE) logs -f $(DB_SERVICE)

.PHONY: db-shell
db-shell: ## Open PostgreSQL shell
	@$(COMPOSE) exec $(DB_SERVICE) psql -U kovago -d kovago

.PHONY: db-reset
db-reset: ## Recreate local PostgreSQL database
	@echo "→ removing PostgreSQL data"
	@$(COMPOSE) down -v
	@$(MAKE) db-up
	@echo "✓ PostgreSQL database recreated"

.PHONY: migration
migration: ## Create a new SQL migration: make migration name=create_shops
	@test -n "$(name)" || (echo "usage: make migration name=create_shops"; exit 1)
	@echo "→ creating migration: $(name)"
	@$(GOOSE) -s -dir $(MIGRATIONS) create "$(name)" sql

.PHONY: migrations-validate
migrations-validate: ## Validate database migrations
	@echo "→ validating migrations"
	@$(GOOSE) -dir $(MIGRATIONS) validate
	@echo "✓ migrations valid"

.PHONY: dev
dev: ## Start development environment
	@$(MAKE) db-up
	@$(MAKE) run
