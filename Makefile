.PHONY: build install test lint clean mocks config help

BINARY_NAME=gotext-translate
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=${VERSION}"
GOFLAGS=-trimpath

build:
	go build ${GOFLAGS} ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/gotext-translate

install:
	go install ${GOFLAGS} ${LDFLAGS} ./cmd/gotext-translate

test:
	go test -v ./...

lint:
	golangci-lint run

mocks:
	mockery --dir pkg/translator --name Translator --output pkg/translator/mocks
	mockery --dir pkg/translator --name Provider --output pkg/translator/mocks
	mockery --dir pkg/translator --name Factory --output pkg/translator/mocks

clean:
	rm -rf bin/
	go clean

# Create a sample configuration file
config:
	cp config.example.yaml config.yaml

help:
	@echo "Make targets:"
	@echo "  build    - Build the application"
	@echo "  install  - Install the application"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  mocks    - Generate mock files"
	@echo "  clean    - Remove build artifacts"
	@echo "  config   - Create a config.yaml from the example file"