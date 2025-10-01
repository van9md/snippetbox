.PHONY: build test run fmt

build:
	go build -o bin/app ./...

test:
	go test ./...

run:
	go run ./...

fmt:
	go fmt ./...

