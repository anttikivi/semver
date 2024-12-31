GCI_VERSION = 0.13.5
GOFUMPT_VERSION = 0.7.0
GOLANGCI_LINT_VERSION = 1.62.2

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. ./...

## Formatting tasks

.PHONY: fmt
fmt:
	go mod tidy
	go run github.com/daixiang0/gci@v$(GCI_VERSION) write . --skip-generated --skip-vendor -s standard -s default
	go run mvdan.cc/gofumpt@v$(GOFUMPT_VERSION) -l -w .

## Lint tasks

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION) run ./...
