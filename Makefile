APP := kovago
CMD := ./cmd/kovago
BIN := ./bin/${APP}

GO := go

CONFIG ?= config/local.toml

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
check: fmt vet test ## Run local quality checks
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
