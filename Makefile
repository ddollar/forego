
# the product we're building
NAME := forego
# the product's main package
MAIN := ./src
# fix our gopath
GOPATH := $(GOPATH):$(PWD)

# build and packaging
TARGETS := $(PWD)/bin
PRODUCT := $(TARGETS)/$(NAME)

# build and packaging for release
GITHASH         := $(shell git log --pretty=format:'%h' -n 1)
VERSION         ?= $(GITHASH)
RELEASE_TARGETS  = $(PWD)/target/$(GOOS)_$(GOARCH)
RELEASE_PRODUCT  = $(NAME)-$(VERSION)
RELEASE_ARCHIVE  = $(RELEASE_PRODUCT)-$(GOOS)-$(GOARCH).tgz
RELEASE_PACKAGE  = $(RELEASE_TARGETS)/$(RELEASE_ARCHIVE)
RELEASE_BINARY   = $(RELEASE_TARGETS)/$(RELEASE_PRODUCT)/bin/$(NAME)

# build and install
PREFIX ?= /usr/local

# sources
SRC = $(shell find src -name \*.go -not -path ./src/vendor -print)

.PHONY: all test clean install release build build_release build_formula

all: build

$(PRODUCT): $(SRC)
	go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $@ $(MAIN)

build: $(PRODUCT) ## Build the product

$(RELEASE_BINARY): $(SRC)
	go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $(RELEASE_BINARY) $(MAIN)

$(RELEASE_PACKAGE): $(RELEASE_BINARY)
	(cd $(RELEASE_TARGETS) && tar -zcf $(RELEASE_ARCHIVE) $(RELEASE_PRODUCT))

build_release: $(RELEASE_PACKAGE)

build_formula: build_release
	$(PWD)/tools/update-formula -v $(VERSION) -o $(PWD)/formula/instaunit.rb $(RELEASE_PACKAGE)

release: test ## Build for all supported architectures
	make build_release GOOS=linux GOARCH=amd64
	make build_release GOOS=freebsd GOARCH=amd64
	make build_formula GOOS=darwin GOARCH=amd64

install: build ## Build and install
	install -m 0755 $(PRODUCT) $(PREFIX)/bin/

test: ## Run tests
	go test ./src/...

clean: ## Delete the built product and any generated files
	rm -rf $(TARGETS)
