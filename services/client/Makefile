name := client-service
version ?= latest
docker_image ?= moorara/$(name)

path := $(shell pwd)
build_dir := $(path)/build
report_dir := $(path)/coverage


clean:
	@ rm -rf $(build_dir) $(report_dir)

docker:
	@ docker build --tag $(docker_image):$(version) .

push:
	@ docker image push $(docker_image):$(version)

save-docker:
	@ docker image save --output ./docker.tar $(docker_image):$(version)

load-docker:
	@ docker image load --input docker.tar

up:
	@ VERSION=$(version) \
	  docker-compose up -d $(name)

down:
	@ VERSION=$(version) \
	  docker-compose down


.PHONY: clean
.PHONY: docker push save-docker load-docker
.PHONY: up down
