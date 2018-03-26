docker_image ?= moorara/react-client
docker_tag ?= latest

path := $(shell pwd)
build_dir := $(path)/public
report_dir := $(path)/coverage


clean:
	@ rm -rf $(build_dir) $(report_dir)

docker:
	@ docker build -t $(docker_image):$(docker_tag) .

up:
	@ VERSION=$(docker_tag) \
	  docker-compose up -d react-client

down:
	@ VERSION=$(docker_tag) \
	  docker-compose down


.PHONY: clean
.PHONY: docker up down
