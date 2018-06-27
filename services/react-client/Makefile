docker_tag ?= latest
docker_image ?= moorara/react-client

path := $(shell pwd)
build_dir := $(path)/public
report_dir := $(path)/coverage
docker_dir := /tmp/workspace


clean:
	@ rm -rf $(build_dir) $(report_dir)

docker:
	@ docker build -t $(docker_image):$(docker_tag) .

save-images:
	@ mkdir -p $(docker_dir)
	@ docker image save -o $(docker_dir)/docker.tar $(docker_image):$(docker_tag)

up:
	@ VERSION=$(docker_tag) \
	  docker-compose up -d react-client

down:
	@ VERSION=$(docker_tag) \
	  docker-compose down


.PHONY: clean
.PHONY: docker save-images
.PHONY: up down
