BIN = forego
SRC = $(shell find . -name '*.go')

.PHONY: all build clean lint release test

all: build

build: $(BIN)

clean:
	rm -f $(BIN)

lint: $(SRC)
	go fmt

release:
	bin/release

test: lint build
	go test -v -race -cover ./...

$(BIN): $(SRC)
	go build -o $@
