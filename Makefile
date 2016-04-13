BIN = forego
SRC = $(shell ls *.go)

.PHONY: all build clean install test lint

all: build

build: $(BIN)

clean:
	rm -f $(BIN)

install: forego
	cp $< ${GOPATH}/bin/

lint: $(SRC)
	go fmt

test: lint build
	go test ./... -cover
	cd eg && ../forego start

$(BIN): $(SRC)
	go build -o $@
