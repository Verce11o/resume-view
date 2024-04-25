.PHONY: test format lint gen-proto build migrate-up migrate-down

test:
	go test -v ./...

format:
	gofmt .

lint:
	golangci-lint run

gen-proto:
	protoc -I protos protos/*.proto --go_out=protos/gen/go --go_opt=paths=source_relative --go-grpc_out=protos/gen/go/ --go-grpc_opt=paths=source_relative

build:
	go build -o resume-view -v cmd/main.go

migrate-up:
	goose -dir ./migrations postgres "host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} sslmode=disable" down

all: test format lint gen-proto build
