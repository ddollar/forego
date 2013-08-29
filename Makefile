BIN = forego
SRC = $(shell ls *.go)

all: build

build: $(BIN)

clean:
	rm -f $(BIN)

install: forego
	cp $< ${GOPATH}/bin/

release: build
	go fmt

$(BIN): $(SRC)
	go build -o $@ $(SRC)
