path := $(shell pwd)


images:
	@ cd ../images && make all

services:
	@ cd ../services/site-service && make docker
	@ cd ../services/sensor-service && make docker
	@ cd ../services/react-client && make docker


pull:
	@ docker-compose pull

up: images services
	@ docker-compose up -d fluentd
	@ docker-compose up -d

down:
	@ docker-compose down

clean:
	@ docker volume ls | grep -oe compose_.*_data | xargs docker volume rm


test-up:
	@ docker-compose up -d fluentd
	@ docker-compose up -d traefik site-service sensor-service

test-integration:
	@ cd tests && ./integration.sh


init-data:
	@ docker run \
	    --network compose_local --link mongo \
	    --volume "$(path)/init/sites.json:/data/sites.json" \
	    mongo \
	    mongoimport --host mongo --port 27017 --db sites --collection sites /data/sites.json
	@ docker run \
	    --network compose_local --link postgres \
	    --volume "$(path)/init/pgpass:/root/.pgpass" \
	    --volume "$(path)/init/sensors.sql:/data/sensors.sql" \
	    postgres \
	    psql --host postgres --port 5432 --username root --dbname sensors --file /data/sensors.sql --echo-all


.PHONY: images services
.PHONY: pull up down clean
.PHONY: test-up test-integration
.PHONY: init-data
