default: help

SHELL:=/bin/bash

.PHONY: help
help:
	@echo build
	@echo testacc
	@echo dev

.PHONY: build
build:
	go install

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: dev
dev:
	[[ -f docker_compose/drivers/postgresql.jar ]] \
		|| ( wget https://jdbc.postgresql.org/download/postgresql-42.6.0.jar -O docker_compose/drivers/postgresql.jar \
		&& chmod a+rwx docker_compose/drivers/postgresql.jar )
	cd docker_compose && docker-compose down
	cd docker_compose && docker-compose up

.PHONY: clean
clean:
	cd docker_compose && docker-compose down
	docker volume rm docker_compose_bambooVolume
	rm docker_compose/drivers/postgresql.jar
