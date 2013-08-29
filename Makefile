BIN = forego
SRC = $(shell ls *.go)

all: build

build: $(BIN)

clean:
	rm -f $(BIN)

install: forego
	cp $< ${GOPATH}/bin/

lint: $(SRC)
	go fmt

dist: forego-linux forego-osx forego-windows

forego-linux: $(SRC)
	env GOOS=linux GOARCH=386 go build -o $@

forego-osx: $(SRC)
	env GOOS=darwin GOARCH=386 go build -o $@

forego-windows: $(SRC)
	env GOOS=windows GOARCH=386 go build -o $@

$(BIN): $(SRC)
	go build -o $@ $(SRC)
