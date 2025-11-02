# Makefile for lang-actor

.DEFAULT_GOAL := help

# Build variables
GO := go
GOFLAGS := #-mod=vendor
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m -l"
TESTFLAGS := -v -count=1 -timeout=30s -race -failfast -shuffle=on -coverprofile=coverage.out
LINTFLAGS := #-v
PACKAGES := $(shell $(GO) list ./... | grep -vE '/tools/|/examples/')

# Declare phony targets
.PHONY: help
help: ## Show this help message
	@echo "G-Graph Makefile"
	@echo "==================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m\033[0m"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Static Analysis
.PHONY: lint
lint: ## Run static analysis tools (golint, go vet, etc.)
	@echo "Running static analysis..."
	@$(GO) vet $(LINTFLAGS) $(PACKAGES)

##@ Testing
.PHONY: test
test: lint ## Run all tests with race detection and comprehensive flags
	@$(GO) test $(TESTFLAGS) $(PACKAGES)

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@$(GO) test -v -bench=. -benchmem -timeout=60s $(PACKAGES)

##@ Documentation
.PHONY: doc doc-serve doc-install
doc: ## Generate static documentation in Go standards format
	@echo "Generating Go documentation..."
	@for pkg in $(PACKAGES); do \
		echo ""; \
		echo "=== Package: $$pkg ==="; \
		$(GO) doc -all $$pkg || true; \
	done
	@echo ""
	@echo "✅ Documentation generated"
	@echo "To browse documentation interactively, run: make doc-serve"

doc-serve: ## Serve documentation locally using pkgsite (install with: make doc-install)
	@if ! command -v pkgsite > /dev/null 2>&1; then \
		echo "❌ pkgsite is not installed. Run 'make doc-install' to install it."; \
		exit 1; \
	fi
	@echo "Starting documentation server at http://localhost:8080..."
	@echo "Press Ctrl+C to stop the server"
	@pkgsite -http=:8080

doc-install: ## Install pkgsite for serving Go documentation locally
	@echo "Installing pkgsite..."
	@$(GO) install golang.org/x/pkgsite/cmd/pkgsite@latest
	@echo "✅ pkgsite installed successfully"
	@echo "Run 'make doc-serve' to start the documentation server"

##@ Cleanup
.PHONY: clean clean-test
clean: clean-test ## Clean all generated files

clean-test: ## Clean test artifacts (coverage files, etc.)
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@echo "✅ Test artifacts cleaned"

##@ Graph Examples
.PHONY: run-conditional-ex run-thread-ex run-helloworld-ex run-loop-ex run-persistence-ex run-velvet-ex run-completion-ex run-all-ex
run-conditional-ex: ## Run the conditional graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/conditional/run.go
run-thread-ex: ## Run the threading graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/conditional_thread/run.go
run-helloworld-ex: ## Run the hello world graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/hello_world/run.go
run-loop-ex: ## Run the loop graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/loop/run.go
run-interrupt-ex: ## Run the interrupt graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/interrupt/run.go
run-persistence-ex: ## Run the persistence graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/loop_persistent/run.go
run-velvet-ex: ## Run the Velvet AIW educational Q&A example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/velvet/run.go
run-completion-ex: ## Run the AIW completion agent example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/completion/run.go
run-all-ex: run-conditional-ex run-thread-ex run-helloworld-ex run-loop-ex run-interrupt-ex run-persistence-ex run-velvet-ex run-completion-ex ## Run all graph examples
	@echo "✅ All examples executed"
