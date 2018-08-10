docker_tag ?= latest
docker_image ?= moorara/site-service
docker_test_image ?= site-service-test

path := $(shell pwd)
docker_dir := /tmp/workspace


clean:
	@ rm -rf coverage .nyc_output

docker:
	@ docker build -t $(docker_image):$(docker_tag) .

docker-test:
	@ docker build \
		--file Dockerfile.test \
		--build-arg version=$(docker_tag) \
		--tag $(docker_test_image) .

save-images:
	@ mkdir -p $(docker_dir)
	@ docker image save -o $(docker_dir)/docker.tar $(docker_image):$(docker_tag)
	@ docker image save -o $(docker_dir)/docker-test.tar $(docker_test_image)

up:
	@ VERSION=$(docker_tag) \
	  docker-compose up -d site-service

down:
	@ VERSION=$(docker_tag) \
	  docker-compose down

test: # docker docker-test
	@ VERSION=$(docker_tag) \
	  docker-compose run site-service-test npm run test && \
		docker-compose down

test-component: # docker docker-test
	@ VERSION=$(docker_tag) \
	  docker-compose run site-service-test npm run test:component && \
		docker container logs site-service | grep '^{' | jq . > component-tests.log && \
		docker-compose down


.PHONY: clean
.PHONY: docker docker-test save-images
.PHONY: up down
.PHONY: test test-component
