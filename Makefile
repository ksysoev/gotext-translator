.PHONY: build install test lint clean

BINARY_NAME=gotext-translate
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/gotext-translate

install:
	go install ${LDFLAGS} ./cmd/gotext-translate

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
	go clean

# Create a sample configuration file
config:
	cp config.example.yaml config.yaml
