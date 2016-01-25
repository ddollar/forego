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
	cd eg && ../forego start -f Procfile.error; test $$? -eq 1
	cd fixtures/port_check && ../../forego start -f Procfile.services
	cd fixtures/port_check && ../../forego start -f Procfile.single -c web=10 web

$(BIN): $(SRC)
	godep go build -o $@
