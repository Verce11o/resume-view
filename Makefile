.PHONY: all test build lint

SERVICES = echo-service resume-view

all: $(SERVICES)
	$(foreach service,$(SERVICES),$(MAKE) -C $(service) all;)

build: $(SERVICES)
	$(foreach service,$(SERVICES),$(MAKE) -C $(service) build;)

lint: $(SERVICES)
	$(foreach service,$(SERVICES),$(MAKE) -C $(service) lint;)

test: $(SERVICES)
	$(foreach service,$(SERVICES),$(MAKE) -C $(service) test;)

compose-build:
	docker-compose build

compose-up:
	docker-compose up -d

compose-stop:
	docker-compose stop

compose-down:
	docker-compose down