.PHONY: all build lint test

SERVICES = echo-service resume-view


build: $(SERVICES)
	$(foreach service,$(SERVICES),$(MAKE) -C $(service) build;)

lint:
	golangci-lint run

test:
	go test -v ./...

gen-proto:
	protoc -I protos protos/*.proto --go_out=protos/gen/go --go_opt=paths=source_relative --go-grpc_out=protos/gen/go/ --go-grpc_opt=paths=source_relative


migrate-up:
	migrate -source file://migrations -database ${MIGRATE_POSTGRES} up

migrate-down:
	migrate -source file://migrations -database ${MIGRATE_POSTGRES} down


compose-build:
	docker-compose build

compose-up:
	docker-compose up -d

compose-stop:
	docker-compose stop

compose-down:
	docker-compose down