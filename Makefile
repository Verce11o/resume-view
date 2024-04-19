include .env

.PHONY: all test build help

test:
	go test ./...

format:
	gofmt ./...

lint:
	golangci-lint run

gen-proto:
	protoc -I protos protos/*.proto --go_out=protos/gen/go --go_opt=paths=source_relative --go-grpc_out=protos/gen/go/ --go-grpc_opt=paths=source_relative

build:
	go build -o resume-view -v cmd/main.go

migrate-up:
	goose -dir ./migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" down

