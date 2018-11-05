images_path := ../../images
services_path := ../../services

images:
	@ cd $(images_path) && make all

services:
	@ cd $(services_path)/react-client && make docker
	@ cd $(services_path)/site-service && make docker
	@ cd $(services_path)/sensor-service && make docker
	@ cd $(services_path)/switch-service && make docker
	@ cd $(services_path)/asset-service && make docker
	@ cd $(services_path)/graphql-server && make docker


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
