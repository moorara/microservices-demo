docker_tag ?= latest
docker_image ?= moorara/sensor-service
docker_test_image ?= sensor-service-test

path := $(shell pwd)
report_dir := $(path)/reports
docker_dir := /tmp/workspace
go_packages := $(shell go list ./... | grep -v //)


clean:
	@ rm -rf $(report_dir)

dep:
	@ dep ensure && \
	  dep ensure -update && \
	  dep prune

run:
	@ go run main.go

build:
	@ go build

docker:
	@ docker build -t $(docker_image):$(docker_tag) .

docker-test:
	@ docker build \
		--file Dockerfile.test \
		--tag $(docker_test_image) .

save-images:
	@ mkdir -p $(docker_dir)
	@ docker image save -o $(docker_dir)/docker.tar $(docker_image):$(docker_tag)
	@ docker image save -o $(docker_dir)/docker-test.tar $(docker_test_image)

up:
	@ VERSION=$(docker_tag) \
	  docker-compose up -d sensor-service

down:
	@ VERSION=$(docker_tag) \
	  docker-compose down

test:
	@ go test -v -race ./...

coverage:
	@ mkdir -p $(report_dir) && \
	  echo "mode: atomic" > $(report_dir)/cover.out
	@ $(foreach package, $(go_packages), \
	    go test -covermode=atomic -coverprofile=cover.out $(package) || exit 1; \
	    tail -n +2 cover.out >> $(report_dir)/cover.out;)
	@ go tool cover -html=$(report_dir)/cover.out -o $(report_dir)/cover.html && \
	  rm cover.out $(report_dir)/cover.out

test-component: # docker up
	@ COMPONENT_TEST=true go test -v ./test

test-component-docker: # docker docker-test
	@ VERSION=$(docker_tag) \
	  docker-compose run sensor-service-test make test-component


.PHONY: clean dep
.PHONY: run build
.PHONY: docker docker-test save-images
.PHONY: up down
.PHONY: test coverage test-component test-component-docker
