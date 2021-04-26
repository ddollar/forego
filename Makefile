BIN = forego
SRC = $(shell find . -name '*.go')

.PHONY: all build clean lint test

all: build

build: $(BIN)

clean:
	rm -f $(BIN)

lint: $(SRC)
	go fmt

test: lint build
	go test -v -race -cover ./...

$(BIN): $(SRC)
	go build -o $@
