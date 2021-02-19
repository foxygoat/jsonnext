# --- Global -------------------------------------------------------------------
O = out

all: build test check-coverage lint  ## build, test, check coverage and lint
	@if [ -e .git/rebase-merge ]; then git --no-pager log -1 --pretty='%h %s'; fi
	@echo '$(COLOUR_GREEN)Success$(COLOUR_NORMAL)'

clean::  ## Remove generated files
	-rm -rf $(O)

.PHONY: all clean

# --- Build --------------------------------------------------------------------

build: | $(O)  ## Build binaries of directories in ./cmd to out/
	go build -o $(O) ./cmd/...
	go build -tags flag -o $(O)/jx-flag ./cmd/jx

install:  ## Build and install binaries in $GOBIN or $GOPATH/bin
	go install ./cmd/...

$(O)/jx: build

.PHONY: build install

# --- Test ---------------------------------------------------------------------
COVERFILE = $(O)/coverage.txt
COVERAGE = 96.5
JSONNET_UNIT = //github.com/yugui/jsonnetunit/raw/master

test: test-go test-jsonnet  ## Run tests and generate a coverage file

test-go: | $(O)
	go test -coverprofile=$(COVERFILE) ./...

test-jsonnet: $(O)/jx
	$(O)/jx -J $(JSONNET_UNIT) lib/jx_test.jsonnet

check-coverage: test  ## Check that test coverage meets the required level
	@go tool cover -func=$(COVERFILE) | $(CHECK_COVERAGE) || $(FAIL_COVERAGE)

cover: test  ## Show test coverage in your browser
	go tool cover -html=$(COVERFILE)

CHECK_COVERAGE = awk -F '[ \t%]+' '/^total:/ {print; if ($$3 < $(COVERAGE)) exit 1}'
FAIL_COVERAGE = { echo '$(COLOUR_RED)FAIL - Coverage below $(COVERAGE)%$(COLOUR_NORMAL)'; exit 1; }

.PHONY: check-coverage cover test

# --- Lint ---------------------------------------------------------------------
GOLINT_VERSION = 1.37.0
GOLINT_INSTALLED_VERSION = $(or $(word 4,$(shell golangci-lint --version 2>/dev/null)),0.0.0)
GOLINT_USE_INSTALLED = $(filter $(GOLINT_INSTALLED_VERSION),v$(GOLINT_VERSION) $(GOLINT_VERSION))
GOLINT = $(if $(GOLINT_USE_INSTALLED),golangci-lint,golangci-lint-v$(GOLINT_VERSION))

GOBIN ?= $(firstword $(subst :, ,$(or $(GOPATH),$(HOME)/go)))/bin

lint: $(if $(GOLINT_USE_INSTALLED),,$(GOBIN)/$(GOLINT))  ## Lint go source code
	$(GOLINT) run

$(GOBIN)/$(GOLINT):
	GOBIN=/tmp go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLINT_VERSION) && \
	      mv /tmp/golangci-lint $@

.PHONY: lint

# --- Utilities ----------------------------------------------------------------
COLOUR_NORMAL = $(shell tput sgr0 2>/dev/null)
COLOUR_RED    = $(shell tput setaf 1 2>/dev/null)
COLOUR_GREEN  = $(shell tput setaf 2 2>/dev/null)
COLOUR_WHITE  = $(shell tput setaf 7 2>/dev/null)

help:
	@awk -F ':.*## ' 'NF == 2 && $$1 ~ /^[A-Za-z0-9_-]+$$/ { printf "$(COLOUR_WHITE)%-30s$(COLOUR_NORMAL)%s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

$(O):
	@mkdir -p $@

.PHONY: help
