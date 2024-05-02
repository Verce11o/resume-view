.PHONY: test format lint build migrate-up migrate-down

test:
	go test -v ./...

format:
	gofmt .

lint:
	golangci-lint run

build:
	go build -o resume-view -v cmd/main.go

all: test format lint build
