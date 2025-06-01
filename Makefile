# This Makefile is POSIX-compliant, and non-compliance is considered a bug. It
# follows the POSIX base specification IEEE Std 1003.1-2024. Documentation can
# be found here: https://pubs.opengroup.org/onlinepubs/9799919799/.

.POSIX:
.SUFFIXES:

ADDLICENSE_VERSION = 1.1.1
GCI_VERSION = 0.13.6
GOFUMPT_VERSION = 0.8.0
GOLANGCI_LINT_VERSION = 2.1.6
GOLINES_VERSION = 0.12.2

# ============================================================================ #
# QUALITY CONTROL
# ============================================================================ #

.PHONY: audit
audit: test lint
	go mod tidy -diff
	go mod verify

.PHONY: lint
lint: install-addlicense install-golangci-lint
	addlicense -check -c "Antti Kivi" -l mit *.go
	golangci-lint run

.PHONY: test
test:
	go test $(GOFLAGS)

.PHONY: bench
bench:
	go test $(GOFLAGS) -bench=.

# ============================================================================ #
# DEVELOPMENT & BUILDING
# ============================================================================ #

.PHONY: tidy
tidy: install-addlicense install-gci install-gofumpt install-golines
	addlicense -c "Antti Kivi" -l mit *.go
	go mod tidy -v
	gci write .
	golines --no-chain-split-dots -w .
	gofumpt -extra -l -w .

.PHONY: fmt
fmt: tidy

.PHONY: fuzz
fuzz:
	@fuzztime="$(FUZZTIME)"; \
	if [ -z "$${fuzztime}" ]; then \
		fuzztime="15s"; \
	fi; \
	echo "Running fuzz tests for $${fuzztime}"; \
	go test $(GOFLAGS) -fuzz="^FuzzParse$$" -fuzztime="$${fuzztime}"; \
	go test $(GOFLAGS) -fuzz=FuzzParseLax -fuzztime="$${fuzztime}"

# ============================================================================ #
# TOOL HELPERS
# ============================================================================ #

.PHONY: install-addlicense
install-addlicense:
	@go install github.com/google/addlicense@v$(ADDLICENSE_VERSION)

.PHONY: install-gci
install-gci:
	@PATH="$${PATH}:$$(go env GOPATH)/bin"; \
	if ! command -v gci >/dev/null 2>&1; then \
		echo "gci not found, installing..."; \
		go install github.com/daixiang0/gci@v$(GCI_VERSION); \
		exit 0; \
	fi; \
	current_version="$$(gci --version 2>/dev/null | awk '{print $$3}')"; \
	if [ "$${current_version}" != "$(GCI_VERSION)" ]; then \
		echo "found gci version $${current_version}, installing version $(GCI_VERSION)..."; \
		go install github.com/daixiang0/gci@v$(GCI_VERSION); \
	fi

.PHONY: install-gofumpt
install-gofumpt:
	@PATH="$${PATH}:$$(go env GOPATH)/bin"; \
	if ! command -v gofumpt >/dev/null 2>&1; then \
		echo "gofumpt not found, installing..."; \
		go install mvdan.cc/gofumpt@v$(GOFUMPT_VERSION); \
		exit 0; \
	fi; \
	current_version="$$(gofumpt --version | awk '{print $$1}' | cut -c 2-)"; \
	if [ "$${current_version}" != "$(GOFUMPT_VERSION)" ]; then \
		echo "found gofumpt version $${current_version}, installing version $(GOFUMPT_VERSION)..."; \
		go install mvdan.cc/gofumpt@v$(GOFUMPT_VERSION); \
	fi

.PHONY: install-golangci-lint
install-golangci-lint:
	@GOPATH="$$(go env GOPATH)"; \
	PATH="$${PATH}:$${GOPATH}/bin"; \
	if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found, installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "${GOPATH}/bin" v$(GOLANGCI_LINT_VERSION); \
		exit 0; \
	fi; \
	current_version=$$(golangci-lint --version 2>/dev/null | awk '{print $$4}'); \
	if [ "$${current_version}" != "$(GOLANGCI_LINT_VERSION)" ]; then \
		echo "found golangci-lint version $${current_version}, installing version $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "${GOPATH}/bin" v$(GOLANGCI_LINT_VERSION); \
	fi

.PHONY: install-golines
install-golines:
	@PATH="$${PATH}:$$(go env GOPATH)/bin"; \
	if ! command -v golines >/dev/null 2>&1; then \
		echo "golines not found, installing..."; \
		./scripts/install_golines "$(GOLINES_VERSION)"; \
		exit 0; \
	fi; \
	current_version="$$(golines --version | head -1 | awk '{print $$2}' | cut -c 2-)"; \
	if [ "$${current_version}" != "$(GOLINES_VERSION)" ]; then \
		echo "found golines version $${current_version}, installing version $(GOLINES_VERSION)..."; \
		./scripts/install_golines "$(GOLINES_VERSION)"; \
	fi
