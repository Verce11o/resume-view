DB_USER=postgres
DB_PASSWORD=vercello
DB_NAME=views
DB_HOST=localhost
DB_PORT=5432


.PHONY: all test build

test:
	go test ./...

build:
	go build -o resume-view -v cmd/main.go

lint:
	golangci-lint run

migrate-up:
	goose -dir ./migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" down