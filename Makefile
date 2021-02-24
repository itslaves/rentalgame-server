.PHONY: all setup clean build test lint

all: lint test build

setup:
	@pre-commit install

clean:
	@rm -rf dist/

build:
	@go build -o dist/rg-server 

test:
	@go test -v ./...

lint:
	@golangci-lint run ./...

swagger:
	@swag init --output internal/docs