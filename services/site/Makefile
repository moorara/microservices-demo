version ?= latest
docker_image ?= moorara/site-service
docker_test_image ?= site-service-test

docker_dir := /tmp/workspace


clean:
	@ rm -rf coverage .nyc_output component-tests.log

docker:
	@ docker build --tag $(docker_image):$(version) .

docker-test:
	@ docker build \
		--file Dockerfile.test \
		--build-arg version=$(version) \
		--tag $(docker_test_image) \
		.

up:
	@ VERSION=$(version) \
	  docker-compose up -d site-service

down:
	@ VERSION=$(version) \
	  docker-compose down

test-docker: # docker docker-test
	@ VERSION=$(version) \
	  docker-compose run test-runner npm run test

test-component-docker: # docker docker-test
	@ VERSION=$(version) \
	  docker-compose run test-runner npm run test:component && \
	  docker container logs site-service | grep '^{' | jq . > component-tests.log && \
	  docker-compose down

push:
	@ docker push $(docker_image):$(version)

save-docker:
	@ docker image save --output ./docker.tar $(docker_image):$(version)

load-docker:
	@ docker image load --input docker.tar


.PHONY: clean
.PHONY: docker docker-test
.PHONY: up down
.PHONY: test-docker test-component-docker
.PHONY: push save-docker load-docker
