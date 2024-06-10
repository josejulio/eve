BINDIR      := $(CURDIR)/bin
BINNAME     ?= eve

BINNAME_CLI     ?= eve_cli

# go option
PKG         := ./...
TAGS        :=
TESTS       := .
TESTFLAGS   :=
LDFLAGS     := -w -s
GOFLAGS     :=
CGO_ENABLED ?= 0

# Rebuild the binary if any of these files change
SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(BINNAME) $(BINDIR)/$(BINNAME_CLI)

run: build
	OPENAI_API_KEY=1 $(BINDIR)/$(BINNAME)

run-cli: build
	$(BINDIR)/$(BINNAME_CLI)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd/eve

$(BINDIR)/$(BINNAME_CLI): $(SRC)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME_CLI) ./cmd/eve_cli
