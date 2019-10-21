name := site-service
version ?= latest
docker_image ?= moorara/$(name)
docker_test_image ?= $(name)-test


clean:
	@ rm -rf coverage .nyc_output component-tests.log

docker:
	@ docker image build --tag $(docker_image):$(version) .

docker-test:
	@ docker image build \
		--file Dockerfile.test \
		--build-arg version=$(version) \
		--tag $(docker_test_image) \
		.

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

test-docker: # docker docker-test
	@ VERSION=$(version) \
	  docker-compose run test-runner npm run test

test-component-docker: # docker docker-test
	@ VERSION=$(version) \
	  docker-compose run test-runner npm run test:component && \
	  docker container logs $(name) | grep '^{' | jq . > component-tests.log && \
	  docker-compose down


.PHONY: clean
.PHONY: docker docker-test
.PHONY: push save-docker load-docker
.PHONY: up down test-docker test-component-docker
