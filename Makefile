.PHONY: all test build help

all:
	+$(MAKE) -C echo-service all
	+$(MAKE) -C resume-view all

compose-build:
	docker-compose build

compose-up:
	docker-compose up -d

compose-stop:
	docker-compose stop

compose-down:
	docker-compose down