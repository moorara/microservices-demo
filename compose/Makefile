path := $(shell pwd)


images:
	@ cd ../images && make all

services:
	@ cd ../services/react-client && make docker
	@ cd ../services/site-service && make docker
	@ cd ../services/sensor-service && make docker
	@ cd ../services/switch-service && make docker
	@ cd ../services/graphql-server && make docker


up: images services
	@ docker-compose up -d fluentd
	@ docker-compose up -d

down:
	@ docker-compose down

clean:
	@ docker volume ls | grep -oe compose_.*_data | xargs docker volume rm


sample-data:
	@ ./seed-data.sh


.PHONY: images services
.PHONY: up down clean
.PHONY: sample-data
