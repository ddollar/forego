
# the product we're building
NAME := forego
# the product's main package
MAIN := ./forego

# fix our gopath
GOPATH := $(GOPATH):$(PWD)
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# build and packaging
TARGETS := $(PWD)/bin
PRODUCT := $(TARGETS)/$(NAME)

# build and packaging for release
GITHASH         := $(shell git log --pretty=format:'%h' -n 1)
BRANCH          := $(shell git rev-parse --abbrev-ref HEAD)
VERSION         ?= $(GITHASH)
RELEASE_TARGETS  = $(PWD)/target/$(GOOS)_$(GOARCH)
RELEASE_PRODUCT  = $(NAME)-$(VERSION)
RELEASE_BUILD    = $(RELEASE_TARGETS)/$(RELEASE_PRODUCT)
RELEASE_BINARY   = $(RELEASE_BUILD)/bin/$(NAME)
RELEASE_ARCHIVE  = $(RELEASE_PRODUCT)-$(GOOS)-$(GOARCH).tgz
RELEASE_PACKAGE  = $(RELEASE_TARGETS)/$(RELEASE_ARCHIVE)

# build and install
PREFIX ?= /usr/local
LATEST ?= latest

# sources
SRC = $(shell find $(MAIN) -name \*.go)

.PHONY: all test clean install build archive

all: build

$(PRODUCT): $(SRC)
	go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $@ $(MAIN)

build: $(PRODUCT) ## Build the product

$(RELEASE_BINARY): $(SRC)
	go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $(RELEASE_BINARY) $(MAIN)

$(RELEASE_PACKAGE): $(RELEASE_BINARY)
	(cd $(RELEASE_TARGETS) && tar -zcf $(RELEASE_ARCHIVE) $(RELEASE_PRODUCT))

archive: $(RELEASE_PACKAGE)

install: build ## Build and install
	@echo "Using sudo to install; you may be prompted for a password..."
	sudo install -m 0755 $(PRODUCT) $(PREFIX)/bin/

test: ## Run tests
	go test $(MAIN)/...

clean: ## Delete the built product and any generated files
	rm -rf $(TARGETS)
