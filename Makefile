BINARY_NAME=go-onedrive-cli

all: build

clean:
	go clean

build: clean
	go build -o $(BINARY_NAME) cmd/goc/*.go

install: build
	cp -f $(BINARY_NAME) $(GOBIN)/