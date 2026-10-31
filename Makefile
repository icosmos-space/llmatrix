.PHONY: help tidy test lint build-proxy fmt vet

MODULES := sdk engine proxy

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

tidy: ## Run go mod tidy in all modules
	@for m in $(MODULES); do \
		echo "==> tidy $$m"; \
		(cd $$m && go mod tidy); \
	done

test: ## Run unit tests in all modules
	@for m in $(MODULES); do \
		echo "==> test $$m"; \
		(cd $$m && go test ./...); \
	done

vet: ## Run go vet in all modules
	@for m in $(MODULES); do \
		echo "==> vet $$m"; \
		(cd $$m && go vet ./...); \
	done

fmt: ## Format all Go sources
	gofmt -w -s .

lint: ## Run golangci-lint
	golangci-lint run ./...

build-proxy: ## Build llmatrix-proxy binary
	cd proxy && go build -o bin/llmatrix-proxy ./cmd/llmatrix-proxy
