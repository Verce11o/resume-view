.PHONY: test format lint build migrate-up migrate-down

test:
	go test -v ./...

format:
	gofmt .

lint:
	golangci-lint run

build:
	go build -o resume-view -v cmd/main.go

migrate-up:
	goose -dir ./migrations postgres "host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} sslmode=disable" down

all: test format lint build
