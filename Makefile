BIN = forego
SRC = $(shell ls *.go)

.PHONY: all build clean lint release test

all: build

build: $(BIN)

clean:
	rm -f $(BIN)

lint: $(SRC)
	go fmt

release:
	@curl -s https://bin.equinox.io/a/gSD5wcgebYp/release-tool-1.8.7-linux-amd64.tar.gz | sudo tar xz -C /usr/local/bin
	@curl -s $(EQUINOX_KEY_URL) -o /tmp/equinox.key
	@equinox release --version=$(shell date +%Y%m%d%H%M%S) --platforms="darwin_386 darwin_amd64 linux_386 linux_amd64 windows_386 windows_amd64" --channel=stable --signing-key=/tmp/equinox.key --app=$(EQUINOX_APP) --token=$(EQUINOX_TOKEN)

test: lint build
	go test -v -race -cover ./...

$(BIN): $(SRC)
	go build -o $@
