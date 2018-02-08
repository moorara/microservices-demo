path := $(shell pwd)
fluentd_version ?= v1.1.0


fluentd:
	@ cd fluentd && \
	  docker build \
	    --build-arg fluentd_version=$(fluentd_version) \
	    -t moorara/fluentd:$(fluentd_version) .

services:
	@ cd ../services/site-service && make docker
	@ cd ../services/sensor-service && make docker


pull:
	@ FLUENTD_VERSION=$(fluentd_version) \
	  docker-compose pull

up: # fluentd services
	@ FLUENTD_VERSION=$(fluentd_version) \
	  docker-compose up -d fluentd
	@ FLUENTD_VERSION=$(fluentd_version) \
	  docker-compose up -d

down:
	@ FLUENTD_VERSION=$(fluentd_version) \
	  docker-compose down

clean:
	@ docker volume ls | grep -oe compose_.*_data | xargs docker volume rm


test-up:
	@ FLUENTD_VERSION=$(fluentd_version) \
	  docker-compose up -d fluentd
	@ FLUENTD_VERSION=$(fluentd_version) \
	  docker-compose up -d traefik site-service sensor-service

test-integration:
	@ cd tests && \
	  ./integration.sh


.PHONY: fluentd services
.PHONY: pull up down clean
.PHONY: test-up test-integration
