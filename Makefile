#
# Generic Makefile for Go programs
#

all: build test check-coverage lint  ## build, test, check coverage and lint
	@if [ -e .git/rebase-merge ]; then git --no-pager log -1 --pretty='%h %s'; fi
	@echo '$(COLOUR_GREEN)Success$(COLOUR_NORMAL)'

clean::  ## Remove generated files

.PHONY: all clean

# Get the first directory in GOPATH
GOPATH1 = $(firstword $(subst :, ,$(GOPATH)))

# -- Build ---------------------------------------------------------------------

BINARIES = $(notdir $(filter-out ./cmd/_%,$(wildcard ./cmd/*)))
INSTALL_DIR = $(or $(GOBIN),$(GOPATH1)/bin,$(HOME)/go/bin)
INSTALLED_BINARIES = $(addprefix $(INSTALL_DIR)/,$(BINARIES))

build: $(BINARIES)  ## Build binaries of directories in ./cmd
install: $(INSTALLED_BINARIES)  ## Build and install binaries in $GOBIN or $GOPATH/bin

ifneq ($(BINARIES),)
$(BINARIES):
	go build -o $@ ./cmd/$@

$(INSTALLED_BINARIES):
	go install ./cmd/$(@F)

clean::
	rm -f $(BINARIES)

# We make $(BINARIES) and $(INSTALLED_BINARIES) .PHONY so we always call the
# Go toolchain - it will decide what needs to be rebuilt.
.PHONY: build $(BINARIES) $(INSTALLED_BINARIES)
endif

# -- Lint ----------------------------------------------------------------------

GOLINT := golangci-lint
GOLINT_DOCKER := docker run --rm -v $(GOPATH1):/go -v $(PWD):/src -w /src golangci/golangci-lint:v1.23.6

# If $(GOLINT) is not installed, use docker to run it
lint: lint-with-$(if $(shell which $(GOLINT)),local,docker)  ## lint the source code

lint-with-local:
	$(GOLINT) run

lint-with-docker:  ## lint the source code with golangci-lint docker image
	$(GOLINT_DOCKER) $(GOLINT) run

.PHONY: lint lint-with-local lint-with-docker

# -- Test ----------------------------------------------------------------------

COVERFILE = coverage.out
COVERAGE = 100

test:  ## Run tests and generate a coverage file
	go test -coverprofile=$(COVERFILE) ./...

check-coverage: test  ## Check that test coverage meets the required level
	@go tool cover -func=$(COVERFILE) | $(CHECK_COVERAGE) || $(FAIL_COVERAGE)

cover: test  ## Show test coverage in your browser
	go tool cover -html=$(COVERFILE)

clean::
	rm -f $(COVERFILE)

CHECK_COVERAGE = awk -F '[ \t%]+' '/^total:/ && $$3 < $(COVERAGE) {exit 1}'
FAIL_COVERAGE = { echo '$(COLOUR_RED)FAIL - Coverage below $(COVERAGE)%$(COLOUR_NORMAL)'; exit 1; }

.PHONY: test check-coverage cover

# --- Utilities ---------------------------------------------------------------

COLOUR_NORMAL = $(shell tput sgr0 2>/dev/null)
COLOUR_RED    = $(shell tput setaf 1 2>/dev/null)
COLOUR_GREEN  = $(shell tput setaf 2 2>/dev/null)
COLOUR_WHITE  = $(shell tput setaf 7 2>/dev/null)

help:
	@awk -F ':.*## ' 'NF == 2 && $$1 ~ /^[A-Za-z0-9_-]+$$/ { printf "$(COLOUR_WHITE)%-30s$(COLOUR_NORMAL)%s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: help
