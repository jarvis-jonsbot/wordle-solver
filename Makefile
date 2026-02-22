.PHONY: build test lint dev

build:
	go build -o bin/ ./...

test:
	go test -race -cover ./...

lint:
	golangci-lint run

dev:
	go run ./cmd/wordle-solver
