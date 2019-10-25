services_path := ../services

services:
	@ cd $(services_path)/client && make docker
	@ cd $(services_path)/site && make docker
	@ cd $(services_path)/sensor && make docker
	@ cd $(services_path)/switch && make docker
	@ cd $(services_path)/asset && make docker
	@ cd $(services_path)/graphql && make docker


up: services
	@ docker-compose up -d fluentd
	@ docker-compose up -d

down:
	@ docker-compose down

clean:
	@ docker volume ls | grep -oe docker_.*_data | xargs docker volume rm


sample-data:
	@ ./seed-data.sh


.PHONY: services
.PHONY: up down clean
.PHONY: sample-data
