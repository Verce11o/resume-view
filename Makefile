.PHONY: all build lint test

SERVICES = echo-service resume-view


build: $(SERVICES)
	$(foreach service,$(SERVICES),$(MAKE) -C $(service) build;)

lint:
	golangci-lint run

test:
	go test -v ./...


compose-build:
	docker-compose build

compose-up:
	docker-compose up -d

compose-stop:
	docker-compose stop

compose-down:
	docker-compose down