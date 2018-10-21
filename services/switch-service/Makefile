docker_tag ?= latest
docker_image ?= moorara/switch-service
docker_test_image ?= switch-service-test

path := $(shell pwd)
docker_dir := /tmp/workspace
coverage_dir := $(path)/coverage


clean:
	@ rm -rf *.log $(coverage_dir)

dep:
	@ dep ensure -update

proto:
	@ ./scripts/proto-gen.sh

run:
	@ go run cmd/main.go

build:
	@ ./scripts/build.sh ./cmd/main.go ./switch-service

docker:
	@ cp -r ../../.git .git
	@ docker build -t $(docker_image):$(docker_tag) .
	@ rm -rf .git

docker-test:
	@ docker build \
		--file Dockerfile.test \
		--tag $(docker_test_image) .

up:
	@ VERSION=$(docker_tag) \
	  docker-compose up -d switch-service

down:
	@ VERSION=$(docker_tag) \
	  docker-compose down

test:
	@ go test -v -race ./...

coverage:
	@ ./scripts/test-unit-cover.sh

test-integration: # docker up
	@ INTEGRATION_TEST=true go test -v ./test/integration_test.go

test-component: # docker up
	@ COMPONENT_TEST=true go test -v ./test/component_test.go

test-integration-docker: # docker docker-test
	@ VERSION=$(docker_tag) \
	  docker-compose run integration-test && \
		docker-compose down

test-component-docker: # docker docker-test
	@ VERSION=$(docker_tag) \
	  docker-compose run component-test && \
		docker container logs switch-service | grep '^{' | jq . > component-tests.log && \
		docker-compose down

push:
	@ docker push $(docker_image):$(docker_tag)

save-images:
	@ mkdir -p $(docker_dir)
	@ docker image save -o $(docker_dir)/docker.tar $(docker_image):$(docker_tag)
	@ docker image save -o $(docker_dir)/docker-test.tar $(docker_test_image)


.PHONY: clean dep
.PHONY: proto run build
.PHONY: docker docker-test
.PHONY: up down
.PHONY: test coverage test-integration test-component test-integration-docker test-component-docker
.PHONY: push save-images
