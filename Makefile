BIN = forego
SRC = $(shell ls *.go)

.PHONY: all build clean install lint release test

all: build

build: $(BIN)

clean:
	rm -f $(BIN)

install: forego
	cp $< ${GOPATH}/bin/

lint: $(SRC)
	go fmt

release:
	# curl -s https://bin.equinox.io/a/gSD5wcgebYp/release-tool-1.8.7-linux-amd64.tar.gz | sudo tar xz -C /usr/local/bin
	echo $(EQUINOX_SIGNING_KEY) > /tmp/equinox.key
	equinox release --version=$(shell date +%s) --channel=stable --signing-key=/tmp/equinox.key --app=$(EQUINOX_APP) --token=$(EQUINOX_TOKEN)


test: lint build
	go test -v -race -cover ./...

$(BIN): $(SRC)
	go build -o $@
