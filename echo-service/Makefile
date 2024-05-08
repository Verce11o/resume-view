.PHONY: test format lint build

test:
	go test -v ./...

format:
	gofmt .

lint:
	golangci-lint run

build:
	go build -o echo-service -v cmd/main.go
	docker build . -t echo-service

all: test format lint build