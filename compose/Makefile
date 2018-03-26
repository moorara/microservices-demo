
images:
	@ cd ../images && make all

services:
	@ cd ../services/site-service && make docker
	@ cd ../services/sensor-service && make docker
	@ cd ../services/react-client && make docker


pull:
	@ FLUENTD_VERSION=$(fluentd_version) \
	  docker-compose pull

up: images services
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


.PHONY: images services
.PHONY: pull up down clean
.PHONY: test-up test-integration
