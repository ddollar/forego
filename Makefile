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

release: forego-darwin-386 forego-darwin-amd64 forego-linux-386 forego-linux-amd64 forego-windows-386.exe forego-windows-amd64.exe

forego-darwin-386: $(SRC)
	env GOOS=darwin GOARCH=386 go build -o $@

forego-darwin-amd64: $(SRC)
	env GOOS=darwin GOARCH=amd64 go build -o $@

forego-linux-386: $(SRC)
	env GOOS=linux GOARCH=386 go build -o $@

forego-linux-amd64: $(SRC)
	env GOOS=linux GOARCH=amd64 go build -o $@

forego-windows-386.exe: $(SRC)
	env GOOS=windows GOARCH=386 go build -o $@

forego-windows-amd64.exe: $(SRC)
	env GOOS=windows GOARCH=amd64 go build -o $@

$(BIN): $(SRC)
	go build -o $@
