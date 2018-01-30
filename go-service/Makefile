docker_image ?= moorara/toys-go-service
docker_tag ?= latest

path := $(shell pwd)
report_dir := $(path)/reports
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

up:
	@ VERSION=$(docker_tag) \
	  docker-compose up -d go-service

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

test-component: up
	@ echo
	@ COMPONENT_TEST=true \
		go test -v ./test


.PHONY: clean dep
.PHONY: run build docker up down
.PHONY: test coverage test-component
