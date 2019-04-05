
# the product we're building
NAME := forego
# the product's main package
MAIN := ./src

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
SRC = $(shell find src -name \*.go -not -path ./src/vendor -print)

.PHONY: all test clean install release build archive publish release formula

all: build

$(PRODUCT): $(SRC)
	go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $@ $(MAIN)

build: $(PRODUCT) ## Build the product

$(RELEASE_BINARY): $(SRC)
	go build -ldflags="-X main.version=$(VERSION) -X main.githash=$(GITHASH)" -o $(RELEASE_BINARY) $(MAIN)

$(RELEASE_PACKAGE): $(RELEASE_BINARY)
	(cd $(RELEASE_TARGETS) && tar -zcf $(RELEASE_ARCHIVE) $(RELEASE_PRODUCT))

archive: $(RELEASE_PACKAGE)

publish: archive
	aws s3 cp --acl public-read $(RELEASE_PACKAGE) s3://bww-artifacts/forego/$(VERSION)/$(RELEASE_ARCHIVE)

formula: archive
	mkdir -p $(RELEASE_BUILD)/formula && $(PWD)/build/update-formula -v $(VERSION) -o $(RELEASE_BUILD)/formula/forego.rb $(RELEASE_PACKAGE)
	aws s3 cp --acl public-read $(RELEASE_BUILD)/formula/forego.rb s3://bww-artifacts/forego/$(LATEST)/forego.rb
	aws s3 cp --acl public-read $(RELEASE_BUILD)/formula/forego.rb s3://bww-artifacts/forego/$(VERSION)/forego.rb

gate:
	@echo && echo "AWS Profile: $(AWS_PROFILE)" && echo "    Version: $(VERSION)" && echo "     Branch: $(BRANCH)"
	@echo && read -p "Release version $(VERSION)? [y/N] " -r continue && echo && [ "$${continue:-N}" = "y" ]

release: gate test ## Build for all supported architectures
	make publish GOOS=linux GOARCH=amd64
	make publish GOOS=freebsd GOARCH=amd64
	make publish formula GOOS=darwin GOARCH=amd64
	@echo && echo "Tag this release:\n\t$ git commit -a -m \"Version $(VERSION)\" && git tag $(VERSION)" && echo

install: build ## Build and install
	install -m 0755 $(PRODUCT) $(PREFIX)/bin/

test: ## Run tests
	go test ./src/...

clean: ## Delete the built product and any generated files
	rm -rf $(TARGETS)
