.PHONY: test format lint build

test:
	go test -v ./...

format:
	gofmt .

lint:
	golangci-lint run

build:
	go build -o employee-service -v cmd/main.go
	docker build . -t employee-service

oapi-gen:
	oapi-codegen -config oapi-config.yml api/api.yml


all: test format lint build